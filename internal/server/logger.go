package server

import (
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/optivor/optivor/internal/config"
)

// NewLogger creates a configured slog.Logger based on config LogLevel and LogFormat settings.
func NewLogger(cfg *config.Config, out io.Writer) *slog.Logger {
	if out == nil {
		out = os.Stdout
	}
	var level slog.Level
	if cfg != nil {
		switch strings.ToLower(cfg.Server.LogLevel) {
		case "debug":
			level = slog.LevelDebug
		case "warn", "warning":
			level = slog.LevelWarn
		case "error":
			level = slog.LevelError
		default:
			level = slog.LevelInfo
		}
	} else {
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	if cfg != nil && strings.ToLower(cfg.Server.LogFormat) == "json" {
		return slog.New(slog.NewJSONHandler(out, opts))
	}
	return slog.New(slog.NewTextHandler(out, opts))
}
