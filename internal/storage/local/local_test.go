package local_test

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/optivor/optivor/internal/storage"
	"github.com/optivor/optivor/internal/storage/local"
)

func TestLocalDriver_Get(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	expectedContent := "hello optivor local storage"

	if err := os.WriteFile(testFile, []byte(expectedContent), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	d, err := local.New(tmpDir)
	if err != nil {
		t.Fatalf("failed to create local driver: %v", err)
	}

	reader, err := d.Get(context.Background(), "test.txt")
	if err != nil {
		t.Fatalf("expected no error getting test.txt, got %v", err)
	}
	defer reader.Close()

	content, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("failed to read content: %v", err)
	}

	if string(content) != expectedContent {
		t.Errorf("expected %q, got %q", expectedContent, string(content))
	}

	// Test non-existent key returns ErrNotFound
	_, err = d.Get(context.Background(), "nonexistent.jpg")
	if err != storage.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}
