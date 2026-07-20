package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunInit(t *testing.T) {
	tempDir := t.TempDir()
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}
	defer func() {
		_ = os.Chdir(origWd)
	}()

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to chdir to tempDir: %v", err)
	}

	tests := []struct {
		name      string
		force     bool
		preCreate bool
		wantErr   bool
	}{
		{
			name:      "Initial creation successful",
			force:     false,
			preCreate: false,
			wantErr:   false,
		},
		{
			name:      "Fails when file exists and force is false",
			force:     false,
			preCreate: true,
			wantErr:   true,
		},
		{
			name:      "Overwrites when file exists and force is true",
			force:     true,
			preCreate: true,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := filepath.Join(tempDir, "optivor.yaml")
			if tt.preCreate {
				_ = os.WriteFile(target, []byte("existing: true"), 0644)
			} else {
				_ = os.Remove(target)
			}

			err := RunInit(tt.force)
			if (err != nil) != tt.wantErr {
				t.Errorf("RunInit(%v) error = %v, wantErr %v", tt.force, err, tt.wantErr)
			}

			if !tt.wantErr {
				if _, err := os.Stat(target); os.IsNotExist(err) {
					t.Errorf("expected optivor.yaml to exist after RunInit")
				}
			}
		})
	}
}
