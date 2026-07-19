package cache

import "github.com/optivor/optivor/internal/pipeline"

type Cache interface {
	Get(key string, params pipeline.TransformParams) (data []byte, contentType string, hit bool, err error)
	Set(key string, params pipeline.TransformParams, data []byte, contentType string) error
}
