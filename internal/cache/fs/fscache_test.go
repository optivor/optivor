package fs_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/optivor/optivor/internal/cache/fs"
	"github.com/optivor/optivor/internal/pipeline"
)

func TestFSCache_SetAndGet(t *testing.T) {
	tmpDir := t.TempDir()
	cacheStore, err := fs.New(tmpDir, 0)
	if err != nil {
		t.Fatalf("failed to create FSCache: %v", err)
	}

	ctx := context.Background()

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
	_, _, hit, err := cacheStore.Get(ctx, key, params)
	if err != nil {
		t.Fatalf("unexpected error on get: %v", err)
	}
	if hit {
		t.Error("expected cache miss, got hit")
	}

	// Set cache entry
	if err := cacheStore.Set(ctx, key, params, content, ct); err != nil {
		t.Fatalf("failed to set cache: %v", err)
	}

	// Second Get should be cache hit
	retData, retCT, hit, err := cacheStore.Get(ctx, key, params)
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

func TestFSCache_LRUEviction(t *testing.T) {
	tmpDir := t.TempDir()
	// Set max size limit to 100 bytes (each entry metadata + data takes ~30 bytes)
	maxSize := int64(100)
	cacheStore, err := fs.New(tmpDir, maxSize)
	if err != nil {
		t.Fatalf("failed to create FSCache: %v", err)
	}

	ctx := context.Background()

	ct := "image/jpeg"
	data := []byte("12345678901234567890") // 20 bytes payload + ~11 bytes header = 31 bytes per entry

	p1 := pipeline.TransformParams{Width: 10}
	p2 := pipeline.TransformParams{Width: 20}
	p3 := pipeline.TransformParams{Width: 30}
	p4 := pipeline.TransformParams{Width: 40}

	// Insert 4 entries (4 * ~31 bytes = ~124 bytes > 100 bytes limit)
	_ = cacheStore.Set(ctx, "img1", p1, data, ct)
	time.Sleep(10 * time.Millisecond)
	_ = cacheStore.Set(ctx, "img2", p2, data, ct)
	time.Sleep(10 * time.Millisecond)
	_ = cacheStore.Set(ctx, "img3", p3, data, ct)
	time.Sleep(10 * time.Millisecond)
	_ = cacheStore.Set(ctx, "img4", p4, data, ct)

	// Oldest entry (img1) should have been evicted
	_, _, hit1, _ := cacheStore.Get(ctx, "img1", p1)
	if hit1 {
		t.Error("expected img1 to be evicted by LRU")
	}

	// Latest entry (img4) must still exist
	_, _, hit4, _ := cacheStore.Get(ctx, "img4", p4)
	if !hit4 {
		t.Error("expected img4 to remain in cache")
	}
}
