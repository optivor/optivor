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
	"github.com/optivor/optivor/internal/storage/router"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Server struct {
	cfg     *config.Config
	driver  storage.StorageDriver
	cache   cache.Cache
	pipe    *pipeline.Pipeline
	bucketRouter router.BucketRouter
	router       chi.Router
	srv          *http.Server
	logger       *slog.Logger
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

func (s *Server) SetBucketRouter(r router.BucketRouter) {
	s.bucketRouter = r
}

func (s *Server) setupRouter() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	//nolint:staticcheck // Chi RealIP middleware used internally behind reverse proxy
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	if s.cfg.Server.RequestTimeout > 0 {
		r.Use(middleware.Timeout(s.cfg.Server.RequestTimeout))
	}
	r.Use(RateLimitMiddleware(s.cfg))
	r.Use(SignedURLMiddleware(s.cfg))
	r.Use(IAMAuthMiddleware(s.cfg))
	r.Use(CrawlerProtectionMiddleware(s.cfg))

	r.Get("/healthz", s.handleHealthz)
	r.Get("/health", s.handleHealthz)
	r.Get("/healtz", s.handleHealthz)
	r.Get("/metrics", MetricsHandler().ServeHTTP)
	r.Get("/fetch", s.handleFetch)
	r.Get("/remote", s.handleFetch)
	r.Get("/preset/{presetName}/*", s.handlePreset)
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

	driverToUse := s.driver
	bucketAlias := "default"
	bucketProvider := "s3"
	bucketPolicy := "public"

	if s.bucketRouter != nil {
		parts := strings.SplitN(key, "/", 2)
		if len(parts) == 2 {
			if drv, alias, err := s.bucketRouter.Resolve(r.Context(), parts[0]); err == nil && drv != nil {
				driverToUse = drv
				key = parts[1]
				bucketAlias = alias
				bucketProvider = s.bucketRouter.Provider(alias)
				bucketPolicy = s.bucketRouter.Policy(alias).String()
			}
		}
	}

	span.SetAttributes(
		attribute.String("bucket.alias", bucketAlias),
		attribute.String("bucket.provider", bucketProvider),
		attribute.String("bucket.policy", bucketPolicy),
	)

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
			// #nosec G705
			_, _ = w.Write(cachedData)
			return
		}
	}
	cacheMissesTotal.Inc()

	// 2. Cache Miss -> Run Pipeline
	transformStart := time.Now()
	data, contentType, err := s.pipe.Run(r.Context(), driverToUse, key, params)
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
	// #nosec G705
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
		case pipeline.FitCover, pipeline.FitContain, pipeline.FitFill, pipeline.FitSmart, pipeline.FitFocal:
			params.Fit = pipeline.FitMode(fitStr)
		default:
			return params, fmt.Errorf("invalid 'fit' mode parameter: %s", fitStr)
		}
	} else {
		params.Fit = pipeline.FitCover
	}

	if focalStr := q.Get("focal"); focalStr != "" {
		parts := strings.Split(focalStr, ",")
		if len(parts) == 2 {
			fx, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
			fy, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
			if err1 == nil && err2 == nil {
				params.FocalX = fx
				params.FocalY = fy
				params.Fit = pipeline.FitFocal
			}
		}
	}

	fmtStr := strings.ToLower(q.Get("format"))
	if fmtStr != "" {
		if fmtStr != "webp" && fmtStr != "avif" && fmtStr != "gif" && fmtStr != "mp4" {
			return params, fmt.Errorf("unsupported 'format' parameter: %s", fmtStr)
		}
		params.Format = fmtStr
	}

	// Dynamic Watermarking & Overlays
	if overlay := q.Get("overlay"); overlay != "" {
		params.Overlay = overlay
	}
	if gravity := q.Get("gravity"); gravity != "" {
		params.Gravity = gravity
	}
	if opacityStr := q.Get("opacity"); opacityStr != "" {
		if op, err := strconv.ParseFloat(opacityStr, 64); err == nil {
			params.Opacity = op
		}
	}
	if scaleStr := q.Get("overlay_scale"); scaleStr != "" {
		if sc, err := strconv.ParseFloat(scaleStr, 64); err == nil {
			params.OverlayScale = sc
		}
	}

	// Image Filters
	if blurStr := q.Get("blur"); blurStr != "" {
		if b, err := strconv.ParseFloat(blurStr, 64); err == nil && b > 0 {
			params.Blur = b
		}
	}
	if gsStr := q.Get("grayscale"); gsStr != "" {
		params.Grayscale = strings.ToLower(gsStr) == "true" || gsStr == "1"
	}
	if pxStr := q.Get("pixelate"); pxStr != "" {
		if px, err := strconv.Atoi(pxStr); err == nil && px > 1 {
			params.Pixelate = px
		}
	}

	return params, nil
}

func (s *Server) handlePreset(w http.ResponseWriter, r *http.Request) {
	presetName := chi.URLParam(r, "presetName")
	if presetName == "" {
		http.Error(w, "preset name is required", http.StatusBadRequest)
		return
	}

	preset, exists := s.cfg.Presets[presetName]
	if !exists {
		http.Error(w, fmt.Sprintf("preset '%s' not found", presetName), http.StatusNotFound)
		return
	}

	params, err := s.parseQueryParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	params = pipeline.ApplyPreset(preset, params)

	key := chi.URLParam(r, "*")
	key = strings.TrimPrefix(key, "/")

	if key == "" {
		http.Error(w, "image key is required", http.StatusBadRequest)
		return
	}

	// Resolve bucket router if configured
	driverToUse := s.driver
	if s.bucketRouter != nil {
		parts := strings.SplitN(key, "/", 2)
		if len(parts) == 2 {
			if drv, _, err := s.bucketRouter.Resolve(r.Context(), parts[0]); err == nil && drv != nil {
				driverToUse = drv
				key = parts[1]
			}
		}
	}

	// 1. Cache Check
	cacheKey := fmt.Sprintf("preset:%s:%s", presetName, key)
	if s.cache != nil {
		cachedData, contentType, hit, cacheErr := s.cache.Get(r.Context(), cacheKey, params)
		if cacheErr == nil && hit {
			w.Header().Set("Content-Type", contentType)
			w.Header().Set("X-Optivor-Cache", "HIT")
			w.WriteHeader(http.StatusOK)
			// #nosec G705
			_, _ = w.Write(cachedData)
			return
		}
	}

	data, contentType, err := s.pipe.Run(r.Context(), driverToUse, key, params)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			http.Error(w, "object not found", http.StatusNotFound)
			return
		}
		s.logger.Error("Preset processing error", "preset", presetName, "key", key, "error", err)
		http.Error(w, "preset processing error", http.StatusInternalServerError)
		return
	}

	if s.cache != nil {
		_ = s.cache.Set(r.Context(), cacheKey, params, data, contentType)
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("X-Optivor-Cache", "MISS")
	w.WriteHeader(http.StatusOK)
	// #nosec G705
	_, _ = w.Write(data)
}
