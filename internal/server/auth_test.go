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

func TestMatchPathPrefix(t *testing.T) {
	tests := []struct {
		name         string
		allowedPaths []string
		targetPath   string
		expected     bool
	}{
		{
			name:         "Empty allowed paths permits everything",
			allowedPaths: nil,
			targetPath:   "user-a/avatar.png",
			expected:     true,
		},
		{
			name:         "Wildcard * permits everything",
			allowedPaths: []string{"*"},
			targetPath:   "user-a/avatar.png",
			expected:     true,
		},
		{
			name:         "Prefix match succeeds for child path",
			allowedPaths: []string{"user-a/*"},
			targetPath:   "user-a/photos/vacation.jpg",
			expected:     true,
		},
		{
			name:         "Prefix match fails for different prefix",
			allowedPaths: []string{"user-a/*"},
			targetPath:   "user-b/photos/vacation.jpg",
			expected:     false,
		},
		{
			name:         "Exact match succeeds",
			allowedPaths: []string{"shared/logo.png"},
			targetPath:   "shared/logo.png",
			expected:     true,
		},
		{
			name:         "Exact match fails on different file",
			allowedPaths: []string{"shared/logo.png"},
			targetPath:   "shared/banner.png",
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MatchPathPrefix(tt.allowedPaths, tt.targetPath)
			if result != tt.expected {
				t.Errorf("expected %v for targetPath %s with allowedPaths %v, got %v", tt.expected, tt.targetPath, tt.allowedPaths, result)
			}
		})
	}
}

func TestValidateIAMAccess(t *testing.T) {
	cfg := &config.Config{
		Auth: config.AuthConfig{
			Roles: []config.RoleConfig{
				{
					Name:         "custom-user-a-editor",
					Description:  "Editor role for User A folder",
					Capabilities: []string{"read", "write"},
					AllowedPaths: []string{"users/user-a/*"},
				},
			},
			APIKeys: []config.APIKeyConfig{
				{
					Key:          "key-user-a",
					Name:         "User A Key",
					Role:         "custom-user-a-editor",
					Buckets:      []string{"tenant-bucket"},
				},
				{
					Key:          "key-path-restricted",
					Name:         "Path Restricted Key",
					Buckets:      []string{"*"},
					Scopes:       []string{"read"},
					AllowedPaths: []string{"public/images/*"},
				},
				{
					Key:     "key-admin",
					Name:    "Admin Key",
					Role:    "admin",
					Buckets: []string{"*"},
				},
			},
		},
	}

	// Test 1: User A Key matching role allowed path
	req1 := httptest.NewRequest(http.MethodGet, "/image/tenant-bucket/users/user-a/avatar.png", nil)
	req1.Header.Set("X-API-Key", "key-user-a")
	if !ValidateIAMAccess(req1, cfg, "tenant-bucket", "users/user-a/avatar.png", "read") {
		t.Errorf("expected key-user-a to be valid for users/user-a/avatar.png")
	}

	// Test 2: User A Key failing for User B folder
	req2 := httptest.NewRequest(http.MethodGet, "/image/tenant-bucket/users/user-b/avatar.png", nil)
	req2.Header.Set("X-API-Key", "key-user-a")
	if ValidateIAMAccess(req2, cfg, "tenant-bucket", "users/user-b/avatar.png", "read") {
		t.Errorf("expected key-user-a to be denied for users/user-b/avatar.png")
	}

	// Test 3: Path restricted key matching allowed path
	req3 := httptest.NewRequest(http.MethodGet, "/image/any-bucket/public/images/logo.png", nil)
	req3.Header.Set("X-API-Key", "key-path-restricted")
	if !ValidateIAMAccess(req3, cfg, "any-bucket", "public/images/logo.png", "read") {
		t.Errorf("expected key-path-restricted to be valid for public/images/logo.png")
	}

	// Test 4: Path restricted key failing for private path
	req4 := httptest.NewRequest(http.MethodGet, "/image/any-bucket/private/secret.png", nil)
	req4.Header.Set("X-API-Key", "key-path-restricted")
	if ValidateIAMAccess(req4, cfg, "any-bucket", "private/secret.png", "read") {
		t.Errorf("expected key-path-restricted to be denied for private/secret.png")
	}

	// Test 5: Admin key with built-in role succeeds everywhere
	req5 := httptest.NewRequest(http.MethodGet, "/image/any-bucket/anything/at/all.png", nil)
	req5.Header.Set("X-API-Key", "key-admin")
	if !ValidateIAMAccess(req5, cfg, "any-bucket", "anything/at/all.png", "write") {
		t.Errorf("expected admin key to have full access")
	}
}

