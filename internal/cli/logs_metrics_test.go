package cli

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRunLogs(t *testing.T) {
	err := RunLogs("10", false)
	if err != nil {
		t.Fatalf("RunLogs failed: %v", err)
	}
}

func TestRunMetrics(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("optivor_requests_total 10"))
	}))
	defer ts.Close()

	err := RunMetrics(ts.URL, false)
	if err != nil {
		t.Fatalf("RunMetrics failed for valid endpoint: %v", err)
	}

	errInvalid := RunMetrics("http://localhost:99999/invalid", false)
	if errInvalid == nil {
		t.Fatalf("expected RunMetrics to fail for invalid endpoint")
	}
}
