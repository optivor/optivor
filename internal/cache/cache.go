package cache

import (
	"context"

	"github.com/optivor/optivor/internal/pipeline"
)

type Cache interface {
	Get(ctx context.Context, key string, params pipeline.TransformParams) (data []byte, contentType string, hit bool, err error)
	Set(ctx context.Context, key string, params pipeline.TransformParams, data []byte, contentType string) error
}

