package server

import (
	"net/http"
	"strings"

	"github.com/optivor/optivor/internal/storage/router"
)

func BucketAuthMiddleware(bucketRouter router.BucketRouter, cfgSecret string, signedURLCheck func(r *http.Request) bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := strings.TrimPrefix(r.URL.Path, "/image/")
			parts := strings.SplitN(path, "/", 2)
			alias := ""
			if len(parts) > 0 {
				alias = parts[0]
			}

			if bucketRouter == nil {
				next.ServeHTTP(w, r)
				return
			}

			policy := bucketRouter.Policy(alias)
			switch policy {
			case router.PolicyPrivate:
				http.Error(w, "Forbidden: private bucket", http.StatusForbidden)
				return
			case router.PolicySigned:
				if !signedURLCheck(r) {
					http.Error(w, "Unauthorized: signed URL required for this bucket", http.StatusUnauthorized)
					return
				}
			case router.PolicyPublic:
				fallthrough
			default:
				// Pass through
			}

			next.ServeHTTP(w, r)
		})
	}
}
