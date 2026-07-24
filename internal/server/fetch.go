package server

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/optivor/optivor/internal/pipeline"
)

func (s *Server) handleFetch(w http.ResponseWriter, r *http.Request) {
	if s.cfg.Remote.Enabled == false {
		http.Error(w, "remote image fetching is disabled", http.StatusForbidden)
		return
	}

	targetURLStr := r.URL.Query().Get("url")
	if targetURLStr == "" {
		http.Error(w, "missing 'url' query parameter", http.StatusBadRequest)
		return
	}

	targetURL, err := url.Parse(targetURLStr)
	if err != nil || (targetURL.Scheme != "http" && targetURL.Scheme != "https") {
		http.Error(w, "invalid or unsupported 'url' parameter scheme", http.StatusBadRequest)
		return
	}

	hostname := targetURL.Hostname()
	if hostname == "" {
		http.Error(w, "invalid host in 'url'", http.StatusBadRequest)
		return
	}

	if !isDomainAllowed(hostname, s.cfg.Remote.AllowedDomains) {
		s.logger.Warn("Remote fetch domain rejected", "domain", hostname)
		http.Error(w, fmt.Sprintf("domain '%s' is not whitelisted", hostname), http.StatusForbidden)
		return
	}

	// SSRF Protection: IP Resolution & Validation
	ips, err := net.LookupIP(hostname)
	if err != nil || len(ips) == 0 {
		http.Error(w, "failed to resolve host IP", http.StatusBadRequest)
		return
	}
	for _, ip := range ips {
		if isPrivateIP(ip) {
			s.logger.Warn("SSRF attempt blocked", "domain", hostname, "ip", ip.String())
			http.Error(w, "access to internal/private IP address is prohibited", http.StatusForbidden)
			return
		}
	}

	params, err := s.parseQueryParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Cache Check
	cacheKey := "remote:" + targetURLStr
	if s.cache != nil {
		cachedData, contentType, hit, cacheErr := s.cache.Get(r.Context(), cacheKey, params)
		if cacheErr == nil && hit {
			w.Header().Set("Content-Type", contentType)
			w.Header().Set("X-Optivor-Cache", "HIT")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(cachedData)
			return
		}
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, targetURL.String(), nil)
	if err != nil {
		http.Error(w, "failed to create remote request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("User-Agent", "Optivor/0.9 RemoteFetcher")

	resp, err := client.Do(req)
	if err != nil {
		s.logger.Error("Remote fetch failed", "url", targetURL.String(), "error", err)
		http.Error(w, "failed to fetch remote image", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("remote origin returned status %d", resp.StatusCode), http.StatusBadGateway)
		return
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "failed to read remote image content", http.StatusInternalServerError)
		return
	}

	data, contentType, err := s.pipe.TransformBytes(r.Context(), bodyBytes, params)
	if err != nil {
		if errors.Is(err, pipeline.ErrOversizedImage) {
			http.Error(w, err.Error(), http.StatusRequestEntityTooLarge)
			return
		}
		s.logger.Error("Remote image transform failed", "error", err)
		http.Error(w, "failed to transform remote image", http.StatusInternalServerError)
		return
	}

	if s.cache != nil {
		_ = s.cache.Set(r.Context(), cacheKey, params, data, contentType)
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("X-Optivor-Cache", "MISS")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func isDomainAllowed(hostname string, allowed []string) bool {
	if len(allowed) == 0 {
		return true
	}
	for _, domain := range allowed {
		if domain == "*" {
			return true
		}
		if strings.EqualFold(hostname, domain) || strings.HasSuffix(strings.ToLower(hostname), "."+strings.ToLower(domain)) {
			return true
		}
	}
	return false
}

func isPrivateIP(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsUnspecified() {
		return true
	}

	if ip4 := ip.To4(); ip4 != nil {
		// 10.0.0.0/8
		if ip4[0] == 10 {
			return true
		}
		// 172.16.0.0/12
		if ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31 {
			return true
		}
		// 192.168.0.0/16
		if ip4[0] == 192 && ip4[1] == 168 {
			return true
		}
		// 169.254.0.0/16
		if ip4[0] == 169 && ip4[1] == 254 {
			return true
		}
	}
	return false
}
