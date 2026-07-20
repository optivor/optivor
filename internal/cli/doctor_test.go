package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunDoctor(t *testing.T) {
	tempDir := t.TempDir()

	validConfig := `server:
  port: 8080
storage:
  s3:
    endpoint: "http://localhost:9000"
    bucket: "test-bucket"
cache:
  fs:
    dir: "/tmp/cache"
image:
  contain_background_color: "#ffffff"
`
	validPath := filepath.Join(tempDir, "valid.yaml")
	_ = os.WriteFile(validPath, []byte(validConfig), 0644)

	invalidPath := filepath.Join(tempDir, "nonexistent.yaml")

	tests := []struct {
		name       string
		configPath string
		wantErr    bool
	}{
		{
			name:       "Valid configuration passes doctor check",
			configPath: validPath,
			wantErr:    false,
		},
		{
			name:       "Nonexistent configuration fails doctor check",
			configPath: invalidPath,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RunDoctor(tt.configPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("RunDoctor(%s) error = %v, wantErr %v", tt.configPath, err, tt.wantErr)
			}
		})
	}
}
