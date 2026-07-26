package persistent

import (
	"context"
	"os"
	"testing"

	"github.com/optivor/optivor/internal/pipeline"
)

func TestPersistentCache_SetAndGet(t *testing.T) {
	dir, err := os.MkdirTemp("", "optivor-persistent-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	pc, err := NewPersistentCache(dir)
	if err != nil {
		t.Fatalf("failed to init persistent cache: %v", err)
	}

	ctx := context.Background()
	key := "test-image.jpg"
	params := pipeline.TransformParams{Width: 200, Height: 200, Format: "webp"}
	testData := []byte("fake-image-bytes")
	contentType := "image/webp"

	// 1. Get non-existent
	_, _, hit, err := pc.Get(ctx, key, params)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if hit {
		t.Errorf("expected cache miss, got hit")
	}

	// 2. Set
	if err := pc.Set(ctx, key, params, testData, contentType); err != nil {
		t.Fatalf("failed to set cache: %v", err)
	}

	// 3. Get hit
	data, cType, hit, err := pc.Get(ctx, key, params)
	if err != nil {
		t.Fatalf("unexpected error on get: %v", err)
	}
	if !hit {
		t.Fatalf("expected cache hit, got miss")
	}
	if string(data) != string(testData) {
		t.Errorf("got data %s, expected %s", string(data), string(testData))
	}
	if cType != contentType {
		t.Errorf("got cType %s, expected %s", cType, contentType)
	}
}
