package redis

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/optivor/optivor/internal/pipeline"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
)

type RedisCache struct {
	client *redis.Client
	prefix string
	ttl    time.Duration
}

func New(addr, password string, db int, prefix string, ttl time.Duration) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Ping to ensure connection is working
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	if prefix == "" {
		prefix = "optivor:cache:"
	}

	return &RedisCache{
		client: client,
		prefix: prefix,
		ttl:    ttl,
	}, nil
}

func (r *RedisCache) Close() error {
	return r.client.Close()
}

func (r *RedisCache) generateKey(key string, params pipeline.TransformParams) string {
	raw := fmt.Sprintf("k=%s|w=%d|h=%d|f=%s|fmt=%s|bg=%s",
		key, params.Width, params.Height, params.Fit, params.Format, params.ContainBackgroundColor)
	hash := sha256.Sum256([]byte(raw))
	return r.prefix + hex.EncodeToString(hash[:])
}

func (r *RedisCache) Get(ctx context.Context, key string, params pipeline.TransformParams) ([]byte, string, bool, error) {
	_, span := otel.Tracer("optivor").Start(ctx, "cache.redis.Get")
	defer span.End()

	cacheKey := r.generateKey(key, params)
	val, err := r.client.Get(ctx, cacheKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, "", false, nil
		}
		return nil, "", false, fmt.Errorf("redis get failed: %w", err)
	}

	if len(val) < 1 {
		return nil, "", false, nil
	}

	contentTypeLen := int(val[0])
	if len(val) < 1+contentTypeLen {
		return nil, "", false, nil
	}

	contentType := string(val[1 : 1+contentTypeLen])
	data := val[1+contentTypeLen:]

	return data, contentType, true, nil
}

func (r *RedisCache) Set(ctx context.Context, key string, params pipeline.TransformParams, data []byte, contentType string) error {
	_, span := otel.Tracer("optivor").Start(ctx, "cache.redis.Set")
	defer span.End()

	cacheKey := r.generateKey(key, params)

	if len(contentType) > 255 {
		return fmt.Errorf("content-type string too long for cache metadata")
	}

	buf := make([]byte, 1+len(contentType)+len(data))
	// #nosec G115
	buf[0] = byte(len(contentType))
	copy(buf[1:], []byte(contentType))
	copy(buf[1+len(contentType):], data)

	err := r.client.Set(ctx, cacheKey, buf, r.ttl).Err()
	if err != nil {
		return fmt.Errorf("redis set failed: %w", err)
	}

	return nil
}
