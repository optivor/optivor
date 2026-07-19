package fs_test

import (
	"bytes"
	"testing"

	"github.com/optivor/optivor/internal/cache/fs"
	"github.com/optivor/optivor/internal/pipeline"
)

func TestFSCache_SetAndGet(t *testing.T) {
	tmpDir := t.TempDir()
	cacheStore, err := fs.New(tmpDir)
	if err != nil {
		t.Fatalf("failed to create FSCache: %v", err)
	}

	params := pipeline.TransformParams{
		Width:  100,
		Height: 100,
		Fit:    pipeline.FitCover,
		Format: "webp",
	}

	key := "products/test/image.jpg"
	content := []byte("fake-image-bytes")
	ct := "image/webp"

	// Initial Get should be cache miss
	_, _, hit, err := cacheStore.Get(key, params)
	if err != nil {
		t.Fatalf("unexpected error on get: %v", err)
	}
	if hit {
		t.Error("expected cache miss, got hit")
	}

	// Set cache entry
	if err := cacheStore.Set(key, params, content, ct); err != nil {
		t.Fatalf("failed to set cache: %v", err)
	}

	// Second Get should be cache hit
	retData, retCT, hit, err := cacheStore.Get(key, params)
	if err != nil {
		t.Fatalf("unexpected error on get after set: %v", err)
	}
	if !hit {
		t.Error("expected cache hit, got miss")
	}
	if retCT != ct {
		t.Errorf("expected content type %s, got %s", ct, retCT)
	}
	if !bytes.Equal(retData, content) {
		t.Error("cached data mismatch")
	}
}
