package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/optivor/optivor/internal/config"
)

func TestHandlePreset_NotFound(t *testing.T) {
	cfg := &config.Config{
		Presets: map[string]config.PresetConfig{
			"avatar": {Width: 100, Height: 100},
		},
	}

	s := New(cfg, nil, nil, nil, nil)
	req := httptest.NewRequest("GET", "/preset/nonexistent/test.jpg", nil)
	rec := httptest.NewRecorder()

	s.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404 Not Found for non-existent preset, got %d", rec.Code)
	}
}
