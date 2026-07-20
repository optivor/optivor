package systemd

import (
	"testing"
)

func TestSystemdUnitValidation(t *testing.T) {
	content, err := GetUnitContent()
	if err != nil {
		t.Fatalf("failed to read embedded systemd unit content: %v", err)
	}

	if err := ValidateUnit(content); err != nil {
		t.Errorf("systemd unit validation failed: %v", err)
	}
}
