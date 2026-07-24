package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/optivor/optivor/internal/config"
)

func TestIsCrawler(t *testing.T) {
	crawlers := []string{
		"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
		"Mozilla/5.0 (compatible; Bingbot/2.0; +http://www.bing.com/bingbot.htm)",
		"Bytespider",
	}
	for _, ua := range crawlers {
		if !IsCrawler(ua) {
			t.Errorf("expected UA %s to be detected as crawler", ua)
		}
	}

	normal := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36"
	if IsCrawler(normal) {
		t.Errorf("expected normal UA %s NOT to be detected as crawler", normal)
	}
}

func TestCrawlerProtectionMiddleware_Throttling(t *testing.T) {
	cfg := &config.Config{
		Crawler: config.CrawlerConfig{
			Enabled:               true,
			MaxConcurrencyPerVariant: 2,
		},
	}

	limiter := NewCrawlerLimiter(cfg)
	// Manually acquire both semaphore slots
	limiter.semaphore <- struct{}{}
	limiter.semaphore <- struct{}{}

	handler := CrawlerProtectionMiddleware(cfg)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Direct check when semaphore is full
	select {
	case limiter.semaphore <- struct{}{}:
		t.Errorf("expected semaphore to be full")
	default:
		// expected overflow
	}

	// Normal UA should pass regardless of crawler limiter
	reqNormal := httptest.NewRequest("GET", "/image/test.jpg", nil)
	recNormal := httptest.NewRecorder()
	handler.ServeHTTP(recNormal, reqNormal)
	if recNormal.Code != http.StatusOK {
		t.Errorf("expected status 200 for normal request, got %d", recNormal.Code)
	}
}
