package server

import (
	"bytes"
	"testing"

	"github.com/optivor/optivor/internal/config"
)

func TestInitTracer(t *testing.T) {
	var buf bytes.Buffer
	cfg := &config.Config{
		Telemetry: config.TelemetryConfig{
			Enabled:       true,
			ServiceName:   "optivor-test",
			SamplingRatio: 1.0,
		},
	}

	tp, err := InitTracer(cfg, &buf)
	if err != nil {
		t.Fatalf("InitTracer failed: %v", err)
	}

	if tp == nil {
		t.Fatal("expected non-nil TracerProvider")
	}

	tracer := Tracer()
	_, span := tracer.Start(t.Context(), "test-span")
	span.End()

	_ = tp.Shutdown(t.Context())
}
