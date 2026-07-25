package local

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/optivor/optivor/internal/storage"
)

// Driver implements storage.StorageDriver for local filesystem storage.
type Driver struct {
	baseDir string
}

// New creates a new local storage driver rooted at baseDir.
func New(baseDir string) (*Driver, error) {
	if baseDir == "" {
		baseDir = "./storage"
	}
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create local storage dir %s: %w", baseDir, err)
	}
	return &Driver{baseDir: baseDir}, nil
}

// Get opens and returns a local file given an object key safely without path traversal.
func (d *Driver) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	absBaseDir, err := filepath.Abs(d.baseDir)
	if err != nil {
		return nil, fmt.Errorf("invalid base dir: %w", err)
	}

	cleanKey := filepath.Clean(key)
	cleanKey = strings.TrimPrefix(cleanKey, "/")

	if strings.Contains(cleanKey, "..") {
		return nil, storage.ErrNotFound
	}

	fullPath := filepath.Join(absBaseDir, cleanKey)
	absPath, err := filepath.Abs(fullPath)
	if err != nil || !strings.HasPrefix(absPath, absBaseDir) {
		return nil, storage.ErrNotFound
	}

	f, err := os.Open(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, storage.ErrNotFound
		}
		return nil, err
	}
	return f, nil
}
