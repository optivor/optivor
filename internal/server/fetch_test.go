package server

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/optivor/optivor/internal/config"
)

func TestDomainAllowed(t *testing.T) {
	allowed := []string{"example.com", "cdn.mysite.org"}
	if !isDomainAllowed("example.com", allowed) {
		t.Errorf("expected example.com to be allowed")
	}
	if !isDomainAllowed("sub.example.com", allowed) {
		t.Errorf("expected sub.example.com to be allowed")
	}
	if isDomainAllowed("malicious.com", allowed) {
		t.Errorf("expected malicious.com to be rejected")
	}

	wildcard := []string{"*"}
	if !isDomainAllowed("anything.com", wildcard) {
		t.Errorf("expected anything.com to be allowed with wildcard")
	}
}

func TestPrivateIP(t *testing.T) {
	privateIPs := []string{"127.0.0.1", "10.0.0.5", "172.20.0.1", "192.168.1.1", "169.254.169.254", "::1"}
	for _, ipStr := range privateIPs {
		ip := net.ParseIP(ipStr)
		if !isPrivateIP(ip) {
			t.Errorf("expected %s to be flagged as private IP", ipStr)
		}
	}

	publicIP := net.ParseIP("8.8.8.8")
	if isPrivateIP(publicIP) {
		t.Errorf("expected 8.8.8.8 to NOT be flagged as private IP")
	}
}

func TestHandleFetch_Disabled(t *testing.T) {
	cfg := &config.Config{
		Remote: config.RemoteConfig{
			Enabled: false,
		},
	}
	s := New(cfg, nil, nil, nil, nil)

	req := httptest.NewRequest("GET", "/fetch?url=https://example.com/test.png", nil)
	rec := httptest.NewRecorder()
	s.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Errorf("expected status 403, got %d", rec.Code)
	}
}

func TestHandleFetch_MissingURL(t *testing.T) {
	cfg := &config.Config{
		Remote: config.RemoteConfig{
			Enabled: true,
		},
	}
	s := New(cfg, nil, nil, nil, nil)

	req := httptest.NewRequest("GET", "/fetch", nil)
	rec := httptest.NewRecorder()
	s.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}
