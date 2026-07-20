package systemd

import (
	"embed"
	"fmt"
	"strings"
)

// UnitTemplate contains the embedded systemd unit configuration template.
//
//go:embed optivor.service
var UnitTemplate embed.FS

// GetUnitContent returns the content of optivor.service unit template.
func GetUnitContent() ([]byte, error) {
	return UnitTemplate.ReadFile("optivor.service")
}

// ValidateUnit checks if the unit file contains required systemd sections and directives.
func ValidateUnit(content []byte) error {
	s := string(content)
	if len(s) == 0 {
		return fmt.Errorf("unit file content is empty")
	}
	requiredDirectives := []string{
		"[Unit]",
		"Description=",
		"[Service]",
		"ExecStart=",
		"[Install]",
		"WantedBy=",
	}
	for _, req := range requiredDirectives {
		if !strings.Contains(s, req) {
			return fmt.Errorf("unit file missing required directive: %s", req)
		}
	}
	return nil
}
