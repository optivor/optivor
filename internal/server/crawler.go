package server

import (
	"net/http"
	"strings"

	"github.com/optivor/optivor/internal/config"
)

var knownCrawlers = []string{
	"googlebot",
	"bingbot",
	"bytespider",
	"baiduspider",
	"yandexbot",
	"facebookexternalhit",
	"twitterbot",
	"rogerbot",
	"linkedinbot",
	"embedly",
	"quora link preview",
	"showyouhaveheart",
	"outbrain",
	"pinterest",
	"slackbot",
	"vkShare",
	"w3c_validator",
}

type CrawlerLimiter struct {
	cfg       *config.Config
	semaphore chan struct{}
}

func NewCrawlerLimiter(cfg *config.Config) *CrawlerLimiter {
	limit := 10
	if cfg != nil && cfg.Crawler.MaxConcurrencyPerVariant > 0 {
		limit = cfg.Crawler.MaxConcurrencyPerVariant
	}
	return &CrawlerLimiter{
		cfg:       cfg,
		semaphore: make(chan struct{}, limit),
	}
}

func IsCrawler(ua string) bool {
	uaLower := strings.ToLower(ua)
	for _, bot := range knownCrawlers {
		if strings.Contains(uaLower, bot) {
			return true
		}
	}
	return false
}

func CrawlerProtectionMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	limiter := NewCrawlerLimiter(cfg)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if cfg != nil && cfg.Crawler.Enabled && IsCrawler(r.UserAgent()) {
				select {
				case limiter.semaphore <- struct{}{}:
					defer func() { <-limiter.semaphore }()
					next.ServeHTTP(w, r)
				default:
					w.Header().Set("Retry-After", "1")
					http.Error(w, "crawler concurrency limit exceeded", http.StatusTooManyRequests)
					return
				}
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}
