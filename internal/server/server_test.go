package server_test

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/optivor/optivor/internal/cache/fs"
	"github.com/optivor/optivor/internal/config"
	"github.com/optivor/optivor/internal/pipeline"
	"github.com/optivor/optivor/internal/server"
	"github.com/optivor/optivor/internal/storage"
)

type mockStorage struct {
	files map[string][]byte
}

func (m *mockStorage) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	data, ok := m.files[key]
	if !ok {
		return nil, storage.ErrNotFound
	}
	return io.NopCloser(bytes.NewReader(data)), nil
}

func createTestJPEG(w, h int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			img.Set(x, y, color.RGBA{R: 50, G: 150, B: 250, A: 255})
		}
	}
	var buf bytes.Buffer
	_ = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 80})
	return buf.Bytes()
}

func setupTestServer(t *testing.T, files map[string][]byte) (*server.Server, http.Handler) {
	t.Helper()
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port: 8080,
			Image: config.ServerImage{
				MaxWidth:  2000,
				MaxHeight: 2000,
			},
		},
		Image: config.ImageConfig{
			ContainBackgroundColor: "#ffffff",
		},
	}

	mockStore := &mockStorage{files: files}
	tmpCacheDir := t.TempDir()
	cacheStore, err := fs.New(tmpCacheDir, 0)
	if err != nil {
		t.Fatalf("failed to init cache: %v", err)
	}

	pipe := pipeline.NewPipeline()
	srv := server.New(cfg, mockStore, cacheStore, pipe, nil)
	return srv, srv.Router()
}

func TestHealthz(t *testing.T) {
	_, handler := setupTestServer(t, nil)
	req := httptest.NewRequest("GET", "/healthz", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", rec.Code)
	}
	if rec.Body.String() != "OK" {
		t.Errorf("expected body 'OK', got '%s'", rec.Body.String())
	}
}

func TestImageRoute_SuccessAndCache(t *testing.T) {
	jpgData := createTestJPEG(300, 300)
	files := map[string][]byte{
		"products/123/main.jpg": jpgData,
	}

	_, handler := setupTestServer(t, files)

	targetURL := "/image/products/123/main.jpg?w=100&h=100&fit=cover&format=webp"

	// Request 1: Cache MISS
	req1 := httptest.NewRequest("GET", targetURL, nil)
	rec1 := httptest.NewRecorder()
	handler.ServeHTTP(rec1, req1)

	if rec1.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec1.Code)
	}
	if ct := rec1.Header().Get("Content-Type"); ct != "image/webp" {
		t.Errorf("expected image/webp, got %s", ct)
	}
	if cacheStatus := rec1.Header().Get("X-Optivor-Cache"); cacheStatus != "MISS" {
		t.Errorf("expected MISS, got %s", cacheStatus)
	}

	// Request 2: Cache HIT
	req2 := httptest.NewRequest("GET", targetURL, nil)
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on cached request, got %d", rec2.Code)
	}
	if cacheStatus := rec2.Header().Get("X-Optivor-Cache"); cacheStatus != "HIT" {
		t.Errorf("expected HIT, got %s", cacheStatus)
	}
}

func TestImageRoute_NotFound(t *testing.T) {
	_, handler := setupTestServer(t, map[string][]byte{})

	req := httptest.NewRequest("GET", "/image/nonexistent.jpg", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404 Not Found, got %d", rec.Code)
	}
}

func TestImageRoute_ValidationErrors(t *testing.T) {
	_, handler := setupTestServer(t, nil)

	testCases := []struct {
		name string
		url  string
	}{
		{"WidthExceedsMax", "/image/img.jpg?w=5000"},
		{"NegativeWidth", "/image/img.jpg?w=-10"},
		{"InvalidFit", "/image/img.jpg?fit=unknown"},
		{"InvalidFormat", "/image/img.jpg?format=invalid"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tc.url, nil)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Errorf("expected 400 Bad Request for %s, got %d", tc.url, rec.Code)
			}
		})
	}
}

func TestImageRoute_OversizedImage(t *testing.T) {
	jpgData := createTestJPEG(300, 300) // 90,000 pixels
	files := map[string][]byte{
		"large.jpg": jpgData,
	}

	cfg := &config.Config{
		Server: config.ServerConfig{
			Port: 8080,
			Image: config.ServerImage{
				MaxWidth:  2000,
				MaxHeight: 2000,
			},
		},
		Image: config.ImageConfig{
			ContainBackgroundColor: "#ffffff",
			MaxPixels:              50000, // 50k max pixels < 90k
		},
	}

	mockStore := &mockStorage{files: files}
	pipe := pipeline.NewPipeline()
	srv := server.New(cfg, mockStore, nil, pipe, nil)

	req := httptest.NewRequest("GET", "/image/large.jpg?w=100&h=100", nil)
	rec := httptest.NewRecorder()
	srv.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("expected 413 StatusRequestEntityTooLarge, got %d", rec.Code)
	}
}

func TestInitTracer_OTLPEndpoint(t *testing.T) {
	cfg := &config.Config{
		Telemetry: config.TelemetryConfig{
			Enabled:       true,
			OTLPEndpoint:  "localhost:4317",
			ServiceName:   "optivor-test",
			SamplingRatio: 1.0,
		},
	}
	tp, err := server.InitTracer(cfg, nil)
	if err != nil {
		t.Fatalf("failed to init tracer with OTLPEndpoint: %v", err)
	}
	if tp == nil {
		t.Error("expected non-nil TracerProvider")
	}
}
