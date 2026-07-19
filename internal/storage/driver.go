package storage

import (
	"context"
	"errors"
	"io"
)

// ErrNotFound is a sentinel error returned when an object does not exist in storage.
var ErrNotFound = errors.New("object not found")

// StorageDriver defines the contract for accessing underlying object storage.
type StorageDriver interface {
	Get(ctx context.Context, key string) (io.ReadCloser, error)
}
