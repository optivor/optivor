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

func TestValidateAPIKey(t *testing.T) {
	cfg := &config.Config{
		Auth: config.AuthConfig{
			APIKeys: []config.APIKeyConfig{
				{
					Key:     "secret-admin-key",
					Name:    "Admin",
					Buckets: []string{"*"},
					Scopes:  []string{"*"},
				},
				{
					Key:     "read-only-key",
					Name:    "Reader",
					Buckets: []string{"avatars"},
					Scopes:  []string{"read"},
				},
			},
		},
	}

	reqAdmin := httptest.NewRequest(http.MethodGet, "/image/avatars/user.jpg", nil)
	reqAdmin.Header.Set("X-API-Key", "secret-admin-key")
	if !ValidateAPIKey(reqAdmin, cfg, "avatars", "read") {
		t.Errorf("expected admin key to be valid for all buckets and scopes")
	}

	reqReader := httptest.NewRequest(http.MethodGet, "/image/avatars/user.jpg", nil)
	reqReader.Header.Set("Authorization", "Bearer read-only-key")
	if !ValidateAPIKey(reqReader, cfg, "avatars", "read") {
		t.Errorf("expected read-only key to be valid for avatars read scope")
	}

	reqInvalidScope := httptest.NewRequest(http.MethodDelete, "/image/avatars/user.jpg", nil)
	reqInvalidScope.Header.Set("X-API-Key", "read-only-key")
	if ValidateAPIKey(reqInvalidScope, cfg, "avatars", "write") {
		t.Errorf("expected read-only key to be invalid for write scope")
	}

	reqInvalidKey := httptest.NewRequest(http.MethodGet, "/image/avatars/user.jpg", nil)
	reqInvalidKey.Header.Set("X-API-Key", "invalid-key")
	if ValidateAPIKey(reqInvalidKey, cfg, "avatars", "read") {
		t.Errorf("expected invalid key to be rejected")
	}
}

func TestGenerateDelegatedSignedURL(t *testing.T) {
	secret := "delegation-secret"
	baseURL := "https://optivor.example.com"
	expires := time.Now().Add(1 * time.Hour)

	delegatedURL := GenerateDelegatedSignedURL(baseURL, "private-bucket", "photo.png", "read", expires, secret)
	u, err := url.Parse(delegatedURL)
	if err != nil {
		t.Fatalf("failed to parse generated URL: %v", err)
	}

	if u.Query().Get("bucket") != "private-bucket" {
		t.Errorf("expected bucket 'private-bucket', got '%s'", u.Query().Get("bucket"))
	}
	if u.Query().Get("key") != "photo.png" {
		t.Errorf("expected key 'photo.png', got '%s'", u.Query().Get("key"))
	}
	if u.Query().Get("sig") == "" {
		t.Errorf("expected non-empty sig parameter")
	}
}
