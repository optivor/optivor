package fs

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/optivor/optivor/internal/pipeline"
	"go.opentelemetry.io/otel"
)

type FSCache struct {
	dir          string
	maxSizeBytes int64
}

func New(dir string, maxSizeBytes int64) (*FSCache, error) {
	if dir == "" {
		return nil, fmt.Errorf("cache directory path cannot be empty")
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}
	return &FSCache{dir: dir, maxSizeBytes: maxSizeBytes}, nil
}

func (c *FSCache) generateKey(key string, params pipeline.TransformParams) string {
	raw := fmt.Sprintf("k=%s|w=%d|h=%d|f=%s|fmt=%s|bg=%s",
		key, params.Width, params.Height, params.Fit, params.Format, params.ContainBackgroundColor)
	hash := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(hash[:])
}

func (c *FSCache) Get(ctx context.Context, key string, params pipeline.TransformParams) ([]byte, string, bool, error) {
	_, span := otel.Tracer("optivor").Start(ctx, "cache.Get")
	defer span.End()

	cacheKey := c.generateKey(key, params)
	filePath := filepath.Join(c.dir, cacheKey)

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, "", false, nil
		}
		return nil, "", false, fmt.Errorf("failed to read cache file: %w", err)
	}

	// Update mod time to track LRU access
	now := time.Now()
	_ = os.Chtimes(filePath, now, now)

	if len(fileData) < 1 {
		return nil, "", false, nil
	}

	contentTypeLen := int(fileData[0])
	if len(fileData) < 1+contentTypeLen {
		return nil, "", false, nil
	}

	contentType := string(fileData[1 : 1+contentTypeLen])
	data := fileData[1+contentTypeLen:]

	return data, contentType, true, nil
}

func (c *FSCache) Set(ctx context.Context, key string, params pipeline.TransformParams, data []byte, contentType string) error {
	_, span := otel.Tracer("optivor").Start(ctx, "cache.Set")
	defer span.End()

	cacheKey := c.generateKey(key, params)
	filePath := filepath.Join(c.dir, cacheKey)

	if len(contentType) > 255 {
		return fmt.Errorf("content-type string too long for cache metadata")
	}

	buf := make([]byte, 1+len(contentType)+len(data))
	buf[0] = byte(len(contentType))
	copy(buf[1:], []byte(contentType))
	copy(buf[1+len(contentType):], data)

	// Write atomically using a temporary file
	tmpFile, err := os.CreateTemp(c.dir, "tmp-cache-*")
	if err != nil {
		return fmt.Errorf("failed to create temp cache file: %w", err)
	}
	tmpName := tmpFile.Name()

	if _, err := tmpFile.Write(buf); err != nil {
		tmpFile.Close()
		os.Remove(tmpName)
		return fmt.Errorf("failed to write to temp cache file: %w", err)
	}
	tmpFile.Close()

	if err := os.Rename(tmpName, filePath); err != nil {
		os.Remove(tmpName)
		return fmt.Errorf("failed to commit cache file: %w", err)
	}

	_ = c.evictIfNeeded()

	return nil
}

type fileEntry struct {
	path    string
	size    int64
	modTime time.Time
}

func (c *FSCache) evictIfNeeded() error {
	if c.maxSizeBytes <= 0 {
		return nil
	}

	entries, err := os.ReadDir(c.dir)
	if err != nil {
		return err
	}

	var files []fileEntry
	var totalSize int64

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		files = append(files, fileEntry{
			path:    filepath.Join(c.dir, entry.Name()),
			size:    info.Size(),
			modTime: info.ModTime(),
		})
		totalSize += info.Size()
	}

	if totalSize <= c.maxSizeBytes {
		return nil
	}

	// Sort oldest modTime first
	sort.Slice(files, func(i, j int) bool {
		return files[i].modTime.Before(files[j].modTime)
	})

	for _, f := range files {
		if totalSize <= c.maxSizeBytes {
			break
		}
		if err := os.Remove(f.path); err == nil {
			totalSize -= f.size
		}
	}

	return nil
}
