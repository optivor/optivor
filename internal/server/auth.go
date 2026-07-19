package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/optivor/optivor/internal/config"
)

// GenerateSignature creates a hex-encoded HMAC-SHA256 signature for a given URL path and query values.
func GenerateSignature(path string, query url.Values, secret string) string {
	q := make(url.Values)
	for k, v := range query {
		if k != "sig" {
			q[k] = v
		}
	}

	sigInput := path
	if len(q) > 0 {
		sigInput += "?" + q.Encode()
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(sigInput))
	return hex.EncodeToString(mac.Sum(nil))
}

// SignedURLMiddleware creates a middleware that validates signed URLs when auth.signed_urls.enabled is true.
func SignedURLMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if cfg == nil || !cfg.Auth.SignedURLs.Enabled || r.URL.Path == "/healthz" || r.URL.Path == "/metrics" {
				next.ServeHTTP(w, r)
				return
			}

			q := r.URL.Query()
			sig := q.Get("sig")
			expiresStr := q.Get("expires")

			if sig == "" || expiresStr == "" {
				http.Error(w, "missing signature or expiration parameter", http.StatusUnauthorized)
				return
			}

			expires, err := strconv.ParseInt(expiresStr, 10, 64)
			if err != nil {
				http.Error(w, "invalid expiration parameter", http.StatusUnauthorized)
				return
			}

			if time.Now().Unix() > expires {
				http.Error(w, "signature expired", http.StatusForbidden)
				return
			}

			expectedSig := GenerateSignature(r.URL.Path, q, cfg.Auth.SignedURLs.Secret)
			if !hmac.Equal([]byte(sig), []byte(expectedSig)) {
				http.Error(w, "invalid signature", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
