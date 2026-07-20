package server

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/optivor/optivor/internal/config"
)

func TestNewLoggerJSON(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			LogLevel:  "info",
			LogFormat: "json",
		},
	}

	var buf bytes.Buffer
	logger := NewLogger(cfg, &buf)
	logger.Info("test json log message", "key", "value")

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("expected valid json log output, got error: %v, raw output: %s", err, buf.String())
	}

	if result["msg"] != "test json log message" {
		t.Errorf("expected msg field 'test json log message', got %v", result["msg"])
	}
	if result["key"] != "value" {
		t.Errorf("expected key field 'value', got %v", result["key"])
	}
}

func TestNewLoggerText(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			LogLevel:  "info",
			LogFormat: "text",
		},
	}

	var buf bytes.Buffer
	logger := NewLogger(cfg, &buf)
	logger.Info("test text log message", "key", "value")

	out := buf.String()
	if !bytes.Contains(buf.Bytes(), []byte("test text log message")) {
		t.Errorf("expected text output to contain message, got: %s", out)
	}
}
