package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/optivor/optivor/internal/config"
)

func TestRateLimitMiddleware(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			RateLimit: config.RateLimitConfig{
				Enabled: true,
				RPS:     2,
				Burst:   3,
			},
		},
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	middleware := RateLimitMiddleware(cfg)
	ts := httptest.NewServer(middleware(handler))
	defer ts.Close()

	client := ts.Client()

	// First 3 requests (within burst = 3) should succeed
	for i := 0; i < 3; i++ {
		resp, err := client.Get(ts.URL + "/image/test.jpg")
		if err != nil {
			t.Fatalf("request %d failed: %v", i+1, err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("request %d: expected 200 OK, got %d", i+1, resp.StatusCode)
		}
		_ = resp.Body.Close()
	}

	// 4th request exceeds burst capability -> 429 Too Many Requests
	resp, err := client.Get(ts.URL + "/image/test.jpg")
	if err != nil {
		t.Fatalf("burst request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("expected 429 Too Many Requests, got %d", resp.StatusCode)
	}
	if retryAfter := resp.Header.Get("Retry-After"); retryAfter == "" {
		t.Error("expected Retry-After header to be present")
	}
}
