package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/optivor/optivor/internal/config"
)

func TestMetricsEndpoint(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port: 8080,
		},
		Storage: config.StorageConfig{
			S3: config.S3Config{
				Endpoint: "http://localhost:9000",
				Bucket:   "test-bucket",
			},
		},
		Cache: config.CacheConfig{
			FS: config.FSCacheConfig{
				Dir: t.TempDir(),
			},
		},
	}

	srv := New(cfg, nil, nil, nil, nil)
	ts := httptest.NewServer(srv.Router())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/metrics")
	if err != nil {
		t.Fatalf("failed to GET /metrics: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK for /metrics, got %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	bodyStr := string(bodyBytes)

	if !strings.Contains(bodyStr, "optivor_cache_hits_total") {
		t.Errorf("expected /metrics output to contain 'optivor_cache_hits_total', got: %s", bodyStr)
	}
}
