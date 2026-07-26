package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/optivor/optivor/internal/config"
)

// GenerateSignature creates a hex-encoded HMAC-SHA256 signature for a given URL path and query values.
// The signature input is constructed as path + "?" + query.Encode() (excluding any existing 'sig' parameter).
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

// MatchPathPrefix evaluates whether targetPath satisfies any allowed prefix pattern.
// Pattern format rules:
// - "*" or empty: matches any path
// - "prefix/*" or "prefix/": matches any key starting with "prefix/"
// - exact match (e.g. "users/avatar.png"): matches targetPath exactly
func MatchPathPrefix(allowedPaths []string, targetPath string) bool {
	if len(allowedPaths) == 0 {
		return true
	}
	targetPath = strings.TrimPrefix(targetPath, "/")
	for _, pattern := range allowedPaths {
		pattern = strings.TrimPrefix(pattern, "/")
		if pattern == "*" || pattern == "" {
			return true
		}
		if strings.HasSuffix(pattern, "*") {
			prefix := strings.TrimSuffix(pattern, "*")
			if strings.HasPrefix(targetPath, prefix) {
				return true
			}
		}
		if targetPath == pattern {
			return true
		}
	}
	return false
}

// ValidateAPIKey checks if the request header contains a valid API key with required scope and bucket access.
func ValidateAPIKey(r *http.Request, cfg *config.Config, bucket string, requiredScope string) bool {
	return ValidateIAMAccess(r, cfg, bucket, "", requiredScope)
}

// ValidateIAMAccess checks if the request contains valid authentication and authorization for the given bucket, object path, and required scope/capability.
func ValidateIAMAccess(r *http.Request, cfg *config.Config, bucket string, objectPath string, requiredScope string) bool {
	if cfg == nil || len(cfg.Auth.APIKeys) == 0 {
		return true
	}

	key := r.Header.Get("X-API-Key")
	if key == "" {
		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			key = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	if key == "" {
		return false
	}

	for _, k := range cfg.Auth.APIKeys {
		if k.Key == key {
			var roleCapabilities []string
			var roleAllowedPaths []string

			if k.Role != "" {
				foundRole := false
				for _, role := range cfg.Auth.Roles {
					if strings.EqualFold(role.Name, k.Role) {
						roleCapabilities = role.Capabilities
						roleAllowedPaths = role.AllowedPaths
						foundRole = true
						break
					}
				}

				if !foundRole {
					switch strings.ToLower(k.Role) {
					case "admin":
						roleCapabilities = []string{"*"}
						roleAllowedPaths = []string{"*"}
					case "editor":
						roleCapabilities = []string{"read", "write"}
						roleAllowedPaths = []string{"*"}
					case "reader-path-only":
						roleCapabilities = []string{"read"}
						roleAllowedPaths = []string{"*"}
					}
				}
			}

			// Validate Bucket Access
			bucketAllowed := false
			if len(k.Buckets) == 0 {
				bucketAllowed = true
			} else {
				for _, b := range k.Buckets {
					if b == "*" || b == bucket {
						bucketAllowed = true
						break
					}
				}
			}
			if !bucketAllowed {
				return false
			}

			// Validate Scope / Capabilities
			scopes := k.Scopes
			if len(scopes) == 0 && len(roleCapabilities) > 0 {
				scopes = roleCapabilities
			}

			scopeAllowed := false
			if requiredScope == "" || len(scopes) == 0 {
				scopeAllowed = true
			} else {
				for _, s := range scopes {
					if s == "*" || strings.EqualFold(s, requiredScope) {
						scopeAllowed = true
						break
					}
				}
			}
			if !scopeAllowed {
				return false
			}

			// Validate Path Prefix Access
			allowedPaths := k.AllowedPaths
			if len(allowedPaths) == 0 && len(roleAllowedPaths) > 0 {
				allowedPaths = roleAllowedPaths
			}

			if objectPath != "" && len(allowedPaths) > 0 {
				if !MatchPathPrefix(allowedPaths, objectPath) {
					return false
				}
			}

			return true
		}
	}

	return false
}

// GenerateDelegatedSignedURL generates a dynamic signed URL delegation token for private bucket access.
func GenerateDelegatedSignedURL(baseURL string, bucket string, key string, scope string, expiresAt time.Time, secret string) string {
	q := url.Values{}
	q.Set("bucket", bucket)
	q.Set("key", key)
	q.Set("scope", scope)
	q.Set("expires", strconv.FormatInt(expiresAt.Unix(), 10))

	sig := GenerateSignature("/auth/delegate", q, secret)
	q.Set("sig", sig)

	return baseURL + "/auth/delegate?" + q.Encode()
}
