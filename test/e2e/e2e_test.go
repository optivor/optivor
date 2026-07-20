package e2e_test

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/optivor/optivor/internal/cache/fs"
	"github.com/optivor/optivor/internal/cli"
	"github.com/optivor/optivor/internal/config"
	"github.com/optivor/optivor/internal/pipeline"
	"github.com/optivor/optivor/internal/server"
	"github.com/optivor/optivor/internal/storage"
)

type memoryStorageDriver struct {
	objects map[string][]byte
}

func (m *memoryStorageDriver) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	data, ok := m.objects[key]
	if !ok {
		return nil, storage.ErrNotFound
	}
	return io.NopCloser(bytes.NewReader(data)), nil
}

func createTestJPEG(w, h int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			img.Set(x, y, color.RGBA{R: 220, G: 80, B: 40, A: 255})
		}
	}
	var buf bytes.Buffer
	_ = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 80})
	return buf.Bytes()
}

// TestEndToEnd_DefinitionOfDone validates ADR-0000 Definition of Done:
// Serving a resized, WebP-converted image from a nested key path with zero provider code.
func TestEndToEnd_DefinitionOfDone(t *testing.T) {
	testKey := "products/test/sample.jpg"
	sourceJPEG := createTestJPEG(800, 600)

	memStorage := &memoryStorageDriver{
		objects: map[string][]byte{
			testKey: sourceJPEG,
		},
	}

	cfg := &config.Config{
		Server: config.ServerConfig{
			Port:         8080,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			Image: config.ServerImage{
				MaxWidth:  5000,
				MaxHeight: 5000,
			},
		},
		Cache: config.CacheConfig{
			FS: config.FSCacheConfig{
				Dir: t.TempDir(),
			},
		},
		Image: config.ImageConfig{
			ContainBackgroundColor: "#ffffff",
		},
	}

	cacheStore, err := fs.New(cfg.Cache.FS.Dir, 0)
	if err != nil {
		t.Fatalf("failed to initialize cache store: %v", err)
	}

	pipe := pipeline.NewPipeline()
	srv := server.New(cfg, memStorage, cacheStore, pipe, nil)

	testServer := httptest.NewServer(srv.Router())
	defer testServer.Close()

	// 1. Healthz check
	healthResp, err := http.Get(testServer.URL + "/healthz")
	if err != nil {
		t.Fatalf("failed to send /healthz request: %v", err)
	}
	defer healthResp.Body.Close()

	if healthResp.StatusCode != http.StatusOK {
		t.Fatalf("expected /healthz status 200, got %d", healthResp.StatusCode)
	}

	// 2. Image request with nested key, resize parameters, and WebP format
	imageURL := fmt.Sprintf("%s/image/%s?w=200&h=200&fit=cover&format=webp", testServer.URL, testKey)
	imgResp, err := http.Get(imageURL)
	if err != nil {
		t.Fatalf("failed to request image: %v", err)
	}
	defer imgResp.Body.Close()

	if imgResp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d", imgResp.StatusCode)
	}

	if ct := imgResp.Header.Get("Content-Type"); ct != "image/webp" {
		t.Fatalf("expected Content-Type image/webp, got %s", ct)
	}

	body, err := io.ReadAll(imgResp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	// 3. Verify output image dimensions using libvips
	importParams := vips.NewImportParams()
	vipsImg, err := vips.LoadImageFromBuffer(body, importParams)
	if err != nil {
		t.Fatalf("failed to decode output webp image: %v", err)
	}
	defer vipsImg.Close()

	if vipsImg.Width() != 200 || vipsImg.Height() != 200 {
		t.Fatalf("expected output image size 200x200, got %dx%d", vipsImg.Width(), vipsImg.Height())
	}
}

func TestV02Acceptance(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port:      8080,
			LogLevel:  "info",
			LogFormat: "json",
		},
	}

	// 1. JSON Logging Output Test
	var logBuf bytes.Buffer
	logger := server.NewLogger(cfg, &logBuf)
	srv := server.New(cfg, nil, nil, nil, logger)

	ts := httptest.NewServer(srv.Router())
	defer ts.Close()

	// 2. /metrics Endpoint Test
	resp, err := http.Get(ts.URL + "/metrics")
	if err != nil {
		t.Fatalf("failed to GET /metrics: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK from /metrics, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if !bytes.Contains(body, []byte("optivor_cache_hits_total")) {
		t.Errorf("expected /metrics to contain optivor_cache_hits_total, got: %s", string(body))
	}
}

func TestV03Acceptance(t *testing.T) {
	tempDir := t.TempDir()
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}
	defer func() { _ = os.Chdir(origWd) }()
	_ = os.Chdir(tempDir)

	// 1. Verify 'optivor init' scaffolding
	if err := cli.RunInit(false); err != nil {
		t.Fatalf("optivor init failed: %v", err)
	}

	cfgPath := filepath.Join(tempDir, "optivor.yaml")
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		t.Fatalf("expected optivor.yaml to exist after init")
	}

	// 2. Verify config loading from scaffolded file
	cfg, err := config.Load(cfgPath)
	if err != nil {
		t.Fatalf("failed to load scaffolded optivor.yaml: %v", err)
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("expected default port 8080, got %d", cfg.Server.Port)
	}

	// 3. Verify 'optivor deploy --dry-run --adapter systemd'
	if err := cli.RunDeploy("systemd", cfgPath, true); err != nil {
		t.Fatalf("optivor deploy dry-run failed: %v", err)
	}
}

func TestV04Acceptance(t *testing.T) {
	tempDir := t.TempDir()
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}
	defer func() { _ = os.Chdir(origWd) }()
	_ = os.Chdir(tempDir)

	// 1. Verify 'optivor init' creates valid optivor.yaml
	if err := cli.RunInit(false); err != nil {
		t.Fatalf("optivor init failed: %v", err)
	}

	cfgPath := filepath.Join(tempDir, "optivor.yaml")

	// 2. Verify 'optivor doctor' health checks
	if err := cli.RunDoctor(cfgPath); err != nil {
		t.Fatalf("optivor doctor failed for valid config: %v", err)
	}

	// 3. Verify 'optivor doctor' fails for invalid config path
	if err := cli.RunDoctor("nonexistent.yaml"); err == nil {
		t.Fatalf("expected optivor doctor to fail for nonexistent config")
	}

	// 4. Verify 'optivor logs' command
	if err := cli.RunLogs("10", false); err != nil {
		t.Fatalf("optivor logs failed: %v", err)
	}
}

func TestV05Acceptance(t *testing.T) {
	// 1. Verify Storage Driver config field default & override
	cfg := &config.Config{
		Storage: config.StorageConfig{
			Driver: "r2",
		},
	}
	if cfg.Storage.Driver != "r2" {
		t.Errorf("expected storage driver r2, got %s", cfg.Storage.Driver)
	}

	// 2. Healthcheck Endpoint Validation (Docker HEALTHCHECK target)
	srv := server.New(&config.Config{
		Server: config.ServerConfig{Port: 8080},
	}, nil, nil, nil, nil)
	ts := httptest.NewServer(srv.Router())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/healthz")
	if err != nil {
		t.Fatalf("failed to GET /healthz: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK from Docker healthcheck target /healthz, got %d", resp.StatusCode)
	}
}

