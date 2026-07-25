package redis

import (
	"context"
	"testing"
	"time"

	"github.com/optivor/optivor/internal/pipeline"
)

func TestRedisCache_SetAndGet(t *testing.T) {
	// Try to connect to a local Redis instance. Skip if it is not running.
	rc, err := New("localhost:6379", "", 0, "optivor:test:", 5*time.Second)
	if err != nil {
		t.Skip("Skipping Redis cache test: local Redis is not running or accessible")
		return
	}
	defer rc.Close()

	ctx := context.Background()
	key := "test-image"
	params := pipeline.TransformParams{
		Width:  100,
		Height: 100,
		Fit:    pipeline.FitCover,
		Format: "webp",
	}
	expectedData := []byte("fake-image-bytes")
	expectedContentType := "image/webp"

	// 1. Get before set (cache miss)
	_, _, hit, err := rc.Get(ctx, key, params)
	if err != nil {
		t.Fatalf("unexpected error on cache get: %v", err)
	}
	if hit {
		t.Fatal("expected cache miss, got hit")
	}

	// 2. Set
	err = rc.Set(ctx, key, params, expectedData, expectedContentType)
	if err != nil {
		t.Fatalf("unexpected error on cache set: %v", err)
	}

	// 3. Get after set (cache hit)
	data, contentType, hit, err := rc.Get(ctx, key, params)
	if err != nil {
		t.Fatalf("unexpected error on cache get: %v", err)
	}
	if !hit {
		t.Fatal("expected cache hit, got miss")
	}
	if string(data) != string(expectedData) {
		t.Errorf("expected data %s, got %s", expectedData, data)
	}
	if contentType != expectedContentType {
		t.Errorf("expected content type %s, got %s", expectedContentType, contentType)
	}
}
