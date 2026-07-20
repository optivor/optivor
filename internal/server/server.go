package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/optivor/optivor/internal/cache"
	"github.com/optivor/optivor/internal/config"
	"github.com/optivor/optivor/internal/pipeline"
	"github.com/optivor/optivor/internal/storage"
	"go.opentelemetry.io/otel/trace"
)

type Server struct {
	cfg     *config.Config
	driver  storage.StorageDriver
	cache   cache.Cache
	pipe    *pipeline.Pipeline
	router  chi.Router
	srv     *http.Server
	logger  *slog.Logger
}

func New(cfg *config.Config, driver storage.StorageDriver, cacheStore cache.Cache, pipe *pipeline.Pipeline, logger *slog.Logger) *Server {
	if logger == nil {
		logger = slog.Default()
	}

	s := &Server{
		cfg:    cfg,
		driver: driver,
		cache:  cacheStore,
		pipe:   pipe,
		logger: logger,
	}

	s.setupRouter()
	return s
}

func (s *Server) setupRouter() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	if s.cfg.Server.RequestTimeout > 0 {
		r.Use(middleware.Timeout(s.cfg.Server.RequestTimeout))
	}
	r.Use(RateLimitMiddleware(s.cfg))
	r.Use(SignedURLMiddleware(s.cfg))

	r.Get("/healthz", s.handleHealthz)
	r.Get("/metrics", MetricsHandler().ServeHTTP)
	r.Get("/image/*", s.handleImage)

	s.router = r

	s.srv = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.cfg.Server.Port),
		Handler:      s.router,
		ReadTimeout:  s.cfg.Server.ReadTimeout,
		WriteTimeout: s.cfg.Server.WriteTimeout,
	}
}

func (s *Server) Router() http.Handler {
	return s.router
}

func (s *Server) ListenAndServe() error {
	s.logger.Info("Starting Optivor HTTP server", "port", s.cfg.Server.Port)
	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down Optivor HTTP server")
	return s.srv.Shutdown(ctx)
}

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

func (s *Server) handleImage(w http.ResponseWriter, r *http.Request) {
	ctx, span := Tracer().Start(r.Context(), "HTTP GET /image", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()
	r = r.WithContext(ctx)

	start := time.Now()
	var statusCode = http.StatusOK
	var formatStr = ""
	var fitStr = "cover"

	defer func() {
		statusStr := strconv.Itoa(statusCode)
		requestsTotal.WithLabelValues(statusStr, formatStr, fitStr).Inc()
		requestDuration.WithLabelValues(statusStr).Observe(time.Since(start).Seconds())
	}()

	key := chi.URLParam(r, "*")
	key = strings.TrimPrefix(key, "/")

	if key == "" {
		statusCode = http.StatusBadRequest
		http.Error(w, "image key is required", statusCode)
		return
	}

	params, err := s.parseQueryParams(r)
	if err != nil {
		statusCode = http.StatusBadRequest
		http.Error(w, err.Error(), statusCode)
		return
	}

	formatStr = params.Format
	fitStr = string(params.Fit)

	// 1. Cache Check
	if s.cache != nil {
		cachedData, contentType, hit, cacheErr := s.cache.Get(r.Context(), key, params)
		if cacheErr != nil {
			s.logger.Warn("Cache get failed", "key", key, "error", cacheErr)
		} else if hit {
			cacheHitsTotal.Inc()
			w.Header().Set("Content-Type", contentType)
			w.Header().Set("X-Optivor-Cache", "HIT")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(cachedData)
			return
		}
	}
	cacheMissesTotal.Inc()

	// 2. Cache Miss -> Run Pipeline
	transformStart := time.Now()
	data, contentType, err := s.pipe.Run(r.Context(), s.driver, key, params)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			statusCode = http.StatusNotFound
			http.Error(w, "object not found", statusCode)
			return
		}
		if errors.Is(err, pipeline.ErrOversizedImage) {
			statusCode = http.StatusRequestEntityTooLarge
			http.Error(w, err.Error(), statusCode)
			return
		}
		if errors.Is(r.Context().Err(), context.DeadlineExceeded) || errors.Is(err, context.DeadlineExceeded) {
			statusCode = http.StatusRequestTimeout
			http.Error(w, "request timeout", statusCode)
			return
		}
		// Distinguish storage gateway errors vs image processing errors
		if strings.Contains(err.Error(), "failed to get object from S3") || strings.Contains(err.Error(), "failed to stat object") {
			s.logger.Error("Storage error", "key", key, "error", err)
			statusCode = http.StatusBadGateway
			http.Error(w, "bad gateway storage error", statusCode)
			return
		}

		s.logger.Error("Pipeline processing error", "key", key, "error", err)
		statusCode = http.StatusInternalServerError
		http.Error(w, "internal processing error", statusCode)
		return
	}
	transformDuration.WithLabelValues(formatStr, fitStr).Observe(time.Since(transformStart).Seconds())

	// 3. Cache Set
	if s.cache != nil {
		if cacheErr := s.cache.Set(r.Context(), key, params, data, contentType); cacheErr != nil {
			s.logger.Warn("Cache set failed", "key", key, "error", cacheErr)
		}
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("X-Optivor-Cache", "MISS")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func (s *Server) parseQueryParams(r *http.Request) (pipeline.TransformParams, error) {
	q := r.URL.Query()
	params := pipeline.TransformParams{
		ContainBackgroundColor: s.cfg.Image.ContainBackgroundColor,
		MaxPixels:              s.cfg.Image.MaxPixels,
	}

	if wStr := q.Get("w"); wStr != "" {
		w, err := strconv.Atoi(wStr)
		if err != nil || w < 0 {
			return params, errors.New("invalid 'w' width parameter")
		}
		if w > s.cfg.Server.Image.MaxWidth {
			return params, fmt.Errorf("'w' parameter exceeds max allowed width (%d)", s.cfg.Server.Image.MaxWidth)
		}
		params.Width = w
	}

	if hStr := q.Get("h"); hStr != "" {
		h, err := strconv.Atoi(hStr)
		if err != nil || h < 0 {
			return params, errors.New("invalid 'h' height parameter")
		}
		if h > s.cfg.Server.Image.MaxHeight {
			return params, fmt.Errorf("'h' parameter exceeds max allowed height (%d)", s.cfg.Server.Image.MaxHeight)
		}
		params.Height = h
	}

	fitStr := strings.ToLower(q.Get("fit"))
	if fitStr != "" {
		switch pipeline.FitMode(fitStr) {
		case pipeline.FitCover, pipeline.FitContain, pipeline.FitFill:
			params.Fit = pipeline.FitMode(fitStr)
		default:
			return params, fmt.Errorf("invalid 'fit' mode parameter: %s", fitStr)
		}
	} else {
		params.Fit = pipeline.FitCover
	}

	fmtStr := strings.ToLower(q.Get("format"))
	if fmtStr != "" {
		if fmtStr != "webp" && fmtStr != "avif" {
			return params, fmt.Errorf("unsupported 'format' parameter: %s", fmtStr)
		}
		params.Format = fmtStr
	}

	return params, nil
}
