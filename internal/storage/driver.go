package storage

import (
	"context"
	"errors"
	"io"
	"time"
)

// ErrNotFound is a sentinel error returned when an object does not exist in storage.
var ErrNotFound = errors.New("object not found")

// ObjectInfo holds metadata information for a storage object.
type ObjectInfo struct {
	Key          string
	Size         int64
	ETag         string
	LastModified time.Time
	ContentType  string
}

// StorageDriver defines the core read contract for accessing underlying object storage.
type StorageDriver interface {
	Get(ctx context.Context, key string) (io.ReadCloser, error)
}

// ExtendedDriver defines optional write, delete, and metadata operations for full bucket management.
type ExtendedDriver interface {
	StorageDriver
	Put(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error
	Delete(ctx context.Context, key string) error
	Stat(ctx context.Context, key string) (ObjectInfo, error)
}
