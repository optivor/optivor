package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/optivor/optivor/internal/config"
)

func TestSignedURLMiddleware(t *testing.T) {
	secret := "my-super-secret-key"

	cfgEnabled := &config.Config{
		Auth: config.AuthConfig{
			SignedURLs: config.SignedURLsConfig{
				Enabled: true,
				Secret:  secret,
			},
		},
	}

	cfgDisabled := &config.Config{
		Auth: config.AuthConfig{
			SignedURLs: config.SignedURLsConfig{
				Enabled: false,
			},
		},
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	now := time.Now().Unix()
	futureExpires := strconv.FormatInt(now+3600, 10)
	pastExpires := strconv.FormatInt(now-3600, 10)

	validQuery := url.Values{"w": []string{"100"}, "expires": []string{futureExpires}}
	validSig := GenerateSignature("/image/test.jpg", validQuery, secret)

	pastQuery := url.Values{"w": []string{"100"}, "expires": []string{pastExpires}}
	expiredSig := GenerateSignature("/image/test.jpg", pastQuery, secret)

	tests := []struct {
		name           string
		cfg            *config.Config
		path           string
		query          url.Values
		sig            string
		expectedStatus int
	}{
		{
			name:           "Disabled auth allows unauthenticated request",
			cfg:            cfgDisabled,
			path:           "/image/test.jpg",
			query:          url.Values{"w": []string{"100"}},
			sig:            "",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Valid signature returns 200 OK",
			cfg:            cfgEnabled,
			path:           "/image/test.jpg",
			query:          validQuery,
			sig:            validSig,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Missing signature returns 401 Unauthorized",
			cfg:            cfgEnabled,
			path:           "/image/test.jpg",
			query:          url.Values{"w": []string{"100"}, "expires": []string{futureExpires}},
			sig:            "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Expired signature returns 403 Forbidden",
			cfg:            cfgEnabled,
			path:           "/image/test.jpg",
			query:          pastQuery,
			sig:            expiredSig,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Invalid signature returns 401 Unauthorized",
			cfg:            cfgEnabled,
			path:           "/image/test.jpg",
			query:          validQuery,
			sig:            "invalid-signature-hex",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Healthz endpoint bypasses auth",
			cfg:            cfgEnabled,
			path:           "/healthz",
			query:          url.Values{},
			sig:            "",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := make(url.Values)
			for k, v := range tt.query {
				q[k] = v
			}
			if tt.sig != "" {
				q.Set("sig", tt.sig)
			}

			reqPath := tt.path
			if len(q) > 0 {
				reqPath += "?" + q.Encode()
			}

			req := httptest.NewRequest(http.MethodGet, reqPath, nil)
			rr := httptest.NewRecorder()

			middleware := SignedURLMiddleware(tt.cfg)
			middleware(handler).ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}
