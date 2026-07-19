package server

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"github.com/optivor/optivor/internal/config"
)

type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiter struct {
	mu      sync.Mutex
	ips     map[string]*ipLimiter
	rps     rate.Limit
	burst   int
	enabled bool
}

func NewRateLimiter(cfg config.RateLimitConfig) *RateLimiter {
	return &RateLimiter{
		ips:     make(map[string]*ipLimiter),
		rps:     rate.Limit(cfg.RPS),
		burst:   cfg.Burst,
		enabled: cfg.Enabled,
	}
}

func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	lim, exists := rl.ips[ip]
	if !exists {
		limiter := rate.NewLimiter(rl.rps, rl.burst)
		rl.ips[ip] = &ipLimiter{limiter: limiter, lastSeen: time.Now()}
		return limiter
	}
	lim.lastSeen = time.Now()
	return lim.limiter
}

func RateLimitMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	if cfg == nil || !cfg.Server.RateLimit.Enabled {
		return func(next http.Handler) http.Handler {
			return next
		}
	}

	rl := NewRateLimiter(cfg.Server.RateLimit)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/healthz" || r.URL.Path == "/metrics" {
				next.ServeHTTP(w, r)
				return
			}

			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr
			}

			limiter := rl.getLimiter(ip)
			if !limiter.Allow() {
				w.Header().Set("Retry-After", "1")
				http.Error(w, "too many requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
