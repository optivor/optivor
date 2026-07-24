package server

import (
	"context"
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

// safeTransport prevents DNS rebinding (TOCTOU) attacks by inspecting IP addresses at dial time.
var safeTransport = &http.Transport{
	DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, fmt.Errorf("invalid network address: %w", err)
		}

		ips, err := net.LookupIP(host)
		if err != nil || len(ips) == 0 {
			return nil, errors.New("failed to resolve target hostname")
		}

		for _, ip := range ips {
			if isPrivateIP(ip) {
				return nil, errors.New("connection to private/internal IP address rejected")
			}
		}

		dialer := &net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}
		return dialer.DialContext(ctx, network, net.JoinHostPort(ips[0].String(), port))
	},
}

// validateAndBuildRemoteURL validates scheme, domain whitelist, and private IPs, returning a clean URL string.
func validateAndBuildRemoteURL(rawURL string, allowedDomains []string) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", errors.New("invalid URL syntax")
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", errors.New("unsupported URL scheme")
	}

	hostname := parsed.Hostname()
	if hostname == "" {
		return "", errors.New("missing host in URL")
	}

	if !isDomainAllowed(hostname, allowedDomains) {
		return "", fmt.Errorf("domain '%s' is not whitelisted", hostname)
	}

	ips, err := net.LookupIP(hostname)
	if err != nil || len(ips) == 0 {
		return "", errors.New("failed to resolve host IP")
	}
	for _, ip := range ips {
		if isPrivateIP(ip) {
			return "", errors.New("access to private IP address is prohibited")
		}
	}

	port := parsed.Port()
	hostPort := hostname
	if port != "" {
		hostPort = net.JoinHostPort(hostname, port)
	}

	cleanURL := fmt.Sprintf("%s://%s%s", parsed.Scheme, hostPort, parsed.RequestURI())
	return cleanURL, nil
}

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

	validatedURL, err := validateAndBuildRemoteURL(targetURLStr, s.cfg.Remote.AllowedDomains)
	if err != nil {
		s.logger.Warn("Remote fetch rejected", "url", targetURLStr, "reason", err)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
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
		Timeout:   10 * time.Second,
		Transport: safeTransport,
	}

	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, validatedURL, nil)
	if err != nil {
		http.Error(w, "failed to create remote request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("User-Agent", "Optivor/0.9 RemoteFetcher")

	resp, err := client.Do(req)
	if err != nil {
		s.logger.Error("Remote fetch failed", "url", validatedURL, "error", err)
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
