package cli

import (
	"testing"
)

func TestRunDeploy(t *testing.T) {
	tests := []struct {
		name       string
		adapter    string
		configPath string
		dryRun     bool
		wantErr    bool
	}{
		{
			name:       "Systemd adapter dry-run successful",
			adapter:    "systemd",
			configPath: "optivor.yaml",
			dryRun:     true,
			wantErr:    false,
		},
		{
			name:       "Unsupported adapter returns error",
			adapter:    "invalid-adapter",
			configPath: "optivor.yaml",
			dryRun:     true,
			wantErr:    true,
		},
		{
			name:       "Systemd adapter live run successful",
			adapter:    "systemd",
			configPath: "optivor.yaml",
			dryRun:     false,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RunDeploy(tt.adapter, tt.configPath, tt.dryRun)
			if (err != nil) != tt.wantErr {
				t.Errorf("RunDeploy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
