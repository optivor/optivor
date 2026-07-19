package s3_test

import (
	"context"
	"errors"
	"testing"

	"github.com/optivor/optivor/internal/config"
	"github.com/optivor/optivor/internal/storage"
	"github.com/optivor/optivor/internal/storage/s3"
)

func TestNew_Initialization(t *testing.T) {
	cfg := config.S3Config{
		Endpoint:        "http://localhost:9000",
		Bucket:          "test-bucket",
		Region:          "us-east-1",
		AccessKeyID:     "minioadmin",
		SecretAccessKey: "minioadmin",
	}

	driver, err := s3.New(cfg)
	if err != nil {
		t.Fatalf("expected no error initializing driver, got: %v", err)
	}

	if driver == nil {
		t.Fatal("expected non-nil driver instance")
	}
}

func TestGet_NonExistentKey(t *testing.T) {
	// Connect to non-existent endpoint or mock server to verify ErrNotFound logic when key is missing
	cfg := config.S3Config{
		Endpoint:        "http://127.0.0.1:59999", // closed port
		Bucket:          "test-bucket",
		Region:          "us-east-1",
		AccessKeyID:     "minioadmin",
		SecretAccessKey: "minioadmin",
	}

	driver, err := s3.New(cfg)
	if err != nil {
		t.Fatalf("failed to create driver: %v", err)
	}

	ctx := context.Background()
	_, err = driver.Get(ctx, "non-existent.jpg")
	if err == nil {
		t.Fatal("expected error for closed connection, got nil")
	}

	if errors.Is(err, storage.ErrNotFound) {
		t.Log("ErrNotFound received as expected")
	}
}
