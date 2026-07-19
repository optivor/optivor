package fs

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/optivor/optivor/internal/pipeline"
)

type FSCache struct {
	dir string
}

func New(dir string) (*FSCache, error) {
	if dir == "" {
		return nil, fmt.Errorf("cache directory path cannot be empty")
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}
	return &FSCache{dir: dir}, nil
}

func (c *FSCache) generateKey(key string, params pipeline.TransformParams) string {
	raw := fmt.Sprintf("k=%s|w=%d|h=%d|f=%s|fmt=%s|bg=%s",
		key, params.Width, params.Height, params.Fit, params.Format, params.ContainBackgroundColor)
	hash := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(hash[:])
}

func (c *FSCache) Get(key string, params pipeline.TransformParams) ([]byte, string, bool, error) {
	cacheKey := c.generateKey(key, params)
	filePath := filepath.Join(c.dir, cacheKey)

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, "", false, nil
		}
		return nil, "", false, fmt.Errorf("failed to read cache file: %w", err)
	}

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

func (c *FSCache) Set(key string, params pipeline.TransformParams, data []byte, contentType string) error {
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

	return nil
}
