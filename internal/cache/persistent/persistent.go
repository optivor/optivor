package persistent

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/optivor/optivor/internal/pipeline"
)

type PersistentCache struct {
	baseDir string
}

func NewPersistentCache(baseDir string) (*PersistentCache, error) {
	if baseDir == "" {
		baseDir = "/tmp/optivor-persistent-cache"
	}
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create persistent cache dir: %w", err)
	}
	return &PersistentCache{baseDir: baseDir}, nil
}

func (p *PersistentCache) generateCacheKey(key string, params pipeline.TransformParams) string {
	h := sha256.New()
	h.Write([]byte(key))
	h.Write([]byte(fmt.Sprintf("%d:%d:%s:%s:%s", params.Width, params.Height, params.Fit, params.Format, params.ContainBackgroundColor)))
	return hex.EncodeToString(h.Sum(nil))
}

func (p *PersistentCache) Get(ctx context.Context, key string, params pipeline.TransformParams) ([]byte, string, bool, error) {
	hash := p.generateCacheKey(key, params)
	dataPath := filepath.Join(p.baseDir, hash+".data")
	metaPath := filepath.Join(p.baseDir, hash+".meta")

	data, err := os.ReadFile(dataPath)
	if err != nil {
		return nil, "", false, nil
	}

	metaBytes, err := os.ReadFile(metaPath)
	if err != nil {
		return nil, "", false, nil
	}

	return data, string(metaBytes), true, nil
}

func (p *PersistentCache) Set(ctx context.Context, key string, params pipeline.TransformParams, data []byte, contentType string) error {
	hash := p.generateCacheKey(key, params)
	dataPath := filepath.Join(p.baseDir, hash+".data")
	metaPath := filepath.Join(p.baseDir, hash+".meta")

	if err := os.WriteFile(dataPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write persistent cache data: %w", err)
	}

	if err := os.WriteFile(metaPath, []byte(contentType), 0644); err != nil {
		return fmt.Errorf("failed to write persistent cache metadata: %w", err)
	}

	return nil
}
