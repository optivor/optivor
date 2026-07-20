package router_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	"github.com/optivor/optivor/internal/storage"
	"github.com/optivor/optivor/internal/storage/router"
)

type mockDriver struct {
	data map[string][]byte
	err  error
}

func (m *mockDriver) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	if m.err != nil {
		return nil, m.err
	}
	v, ok := m.data[key]
	if !ok {
		return nil, storage.ErrNotFound
	}
	return io.NopCloser(bytes.NewReader(v)), nil
}

func TestRouter_AliasAndPolicy(t *testing.T) {
	d1 := &mockDriver{data: map[string][]byte{"img1.jpg": []byte("primary")}}
	d2 := &mockDriver{data: map[string][]byte{"img1.jpg": []byte("fallback")}}

	targets := map[string]router.BucketTarget{
		"primary-images": {
			Alias:    "primary-images",
			Driver:   d1,
			Policy:   router.PolicyPublic,
			Fallback: "backup-images",
			Provider: "s3",
		},
		"backup-images": {
			Alias:    "backup-images",
			Driver:   d2,
			Policy:   router.PolicySigned,
			Provider: "r2",
		},
		"private-images": {
			Alias:    "private-images",
			Driver:   d1,
			Policy:   router.PolicyPrivate,
			Provider: "b2",
		},
	}

	r, err := router.NewDefaultRouter(targets, "primary-images")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if p := r.Policy("primary-images"); p != router.PolicyPublic {
		t.Errorf("expected public, got %v", p)
	}
	if p := r.Policy("backup-images"); p != router.PolicySigned {
		t.Errorf("expected signed, got %v", p)
	}
	if p := r.Policy("private-images"); p != router.PolicyPrivate {
		t.Errorf("expected private, got %v", p)
	}

	driver, alias, err := r.Resolve(context.Background(), "primary-images")
	if err != nil || alias != "primary-images" {
		t.Fatalf("resolve failed: %v", err)
	}

	rc, err := driver.Get(context.Background(), "img1.jpg")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	defer rc.Close()
}

func TestRouter_Failover(t *testing.T) {
	d1 := &mockDriver{err: errors.New("primary down")}
	d2 := &mockDriver{data: map[string][]byte{"img1.jpg": []byte("fallback_ok")}}

	targets := map[string]router.BucketTarget{
		"primary": {
			Alias:    "primary",
			Driver:   d1,
			Policy:   router.PolicyPublic,
			Fallback: "backup",
		},
		"backup": {
			Alias:  "backup",
			Driver: d2,
			Policy: router.PolicyPublic,
		},
	}

	r, err := router.NewDefaultRouter(targets, "primary")
	if err != nil {
		t.Fatalf("failed creating router: %v", err)
	}

	driver, _, err := r.Resolve(context.Background(), "primary")
	if err != nil {
		t.Fatalf("resolve error: %v", err)
	}

	rc, err := driver.Get(context.Background(), "img1.jpg")
	if err != nil {
		t.Fatalf("expected fallback success, got %v", err)
	}
	defer rc.Close()

	buf, _ := io.ReadAll(rc)
	if string(buf) != "fallback_ok" {
		t.Errorf("expected fallback_ok, got %s", string(buf))
	}
}

func TestRouter_CircularFallbackValidation(t *testing.T) {
	targets := map[string]router.BucketTarget{
		"a": {Alias: "a", Fallback: "b"},
		"b": {Alias: "b", Fallback: "a"},
	}
	_, err := router.NewDefaultRouter(targets, "a")
	if !errors.Is(err, router.ErrCircularFallback) {
		t.Errorf("expected ErrCircularFallback, got %v", err)
	}
}
