package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/optivor/optivor/internal/config"
)

func TestLoad_ValidFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "optivor.yaml")
	content := `
server:
  port: 9090
  read_timeout: 15s
  write_timeout: 20s
  image:
    max_width: 3000
    max_height: 3000

storage:
  s3:
    endpoint: "https://s3.us-east-1.amazonaws.com"
    bucket: "test-bucket"
    region: "us-east-1"
    access_key_id: "test-key"
    secret_access_key: "test-secret"

cache:
  fs:
    dir: "/tmp/custom-cache"

image:
  contain_background_color: "#000000"
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if cfg.Server.Port != 9090 {
		t.Errorf("expected port 9090, got %d", cfg.Server.Port)
	}
	if cfg.Server.ReadTimeout != 15*time.Second {
		t.Errorf("expected read_timeout 15s, got %v", cfg.Server.ReadTimeout)
	}
	if cfg.Server.Image.MaxWidth != 3000 {
		t.Errorf("expected max_width 3000, got %d", cfg.Server.Image.MaxWidth)
	}
	if cfg.Storage.S3.Bucket != "test-bucket" {
		t.Errorf("expected bucket test-bucket, got %s", cfg.Storage.S3.Bucket)
	}
	if cfg.Cache.FS.Dir != "/tmp/custom-cache" {
		t.Errorf("expected cache dir /tmp/custom-cache, got %s", cfg.Cache.FS.Dir)
	}
	if cfg.Image.ContainBackgroundColor != "#000000" {
		t.Errorf("expected contain_background_color #000000, got %s", cfg.Image.ContainBackgroundColor)
	}
}

func TestLoad_EnvOverride(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "optivor.yaml")
	content := `
storage:
  s3:
    endpoint: "https://s3.amazonaws.com"
    bucket: "initial-bucket"
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	t.Setenv("OPTIVOR_STORAGE_S3_BUCKET", "env-overridden-bucket")

	cfg, err := config.Load(configPath)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if cfg.Storage.S3.Bucket != "env-overridden-bucket" {
		t.Errorf("expected env override bucket env-overridden-bucket, got %s", cfg.Storage.S3.Bucket)
	}
}

func TestValidate_MissingRequiredFields(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port: 8080,
			Image: config.ServerImage{
				MaxWidth:  5000,
				MaxHeight: 5000,
			},
		},
		Storage: config.StorageConfig{
			S3: config.S3Config{
				Endpoint: "", // missing
				Bucket:   "my-bucket",
			},
		},
		Cache: config.CacheConfig{
			FS: config.FSCacheConfig{
				Dir: "/tmp/cache",
			},
		},
	}

	if err := config.Validate(cfg); err == nil {
		t.Errorf("expected validation error for missing endpoint, got nil")
	}
}
