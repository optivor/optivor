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

// Get opens and returns a local file given an object key.
func (d *Driver) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	cleanKey := filepath.Clean(key)
	cleanKey = strings.TrimPrefix(cleanKey, "/")
	fullPath := filepath.Join(d.baseDir, cleanKey)

	f, err := os.Open(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, storage.ErrNotFound
		}
		return nil, err
	}
	return f, nil
}
