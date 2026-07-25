package redis

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/optivor/optivor/internal/pipeline"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
)

type CircuitBreakerState int

const (
	StateClosed CircuitBreakerState = iota
	StateOpen
	StateHalfOpen
)

type RedisCache struct {
	client              *redis.Client
	prefix              string
	ttl                 time.Duration
	poolSize            int
	minIdleConns        int
	maxFailures         int32
	cooldownDuration    time.Duration
	consecutiveFailures atomic.Int32
	state               atomic.Int32 // 0: Closed, 1: Open, 2: HalfOpen
	lastStateChange     time.Time
	mu                  sync.Mutex
	hits                atomic.Uint64
	misses              atomic.Uint64
}

type Config struct {
	Addr             string
	Password         string
	DB               int
	Prefix           string
	TTL              time.Duration
	PoolSize         int
	MinIdleConns     int
	MaxFailures      int
	CooldownDuration time.Duration
}

func New(addr, password string, db int, prefix string, ttl time.Duration) (*RedisCache, error) {
	return NewWithConfig(Config{
		Addr:             addr,
		Password:         password,
		DB:               db,
		Prefix:           prefix,
		TTL:              ttl,
		PoolSize:         10,
		MinIdleConns:     5,
		MaxFailures:      5,
		CooldownDuration: 30 * time.Second,
	})
}

func NewWithConfig(cfg Config) (*RedisCache, error) {
	if cfg.PoolSize <= 0 {
		cfg.PoolSize = 10
	}
	if cfg.MinIdleConns < 0 {
		cfg.MinIdleConns = 2
	}
	if cfg.MaxFailures <= 0 {
		cfg.MaxFailures = 5
	}
	if cfg.CooldownDuration <= 0 {
		cfg.CooldownDuration = 30 * time.Second
	}

	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	})

	// Ping to ensure connection is working
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	if cfg.Prefix == "" {
		cfg.Prefix = "optivor:cache:"
	}

	rc := &RedisCache{
		client:           client,
		prefix:           cfg.Prefix,
		ttl:              cfg.TTL,
		poolSize:         cfg.PoolSize,
		minIdleConns:     cfg.MinIdleConns,
		maxFailures:      int32(cfg.MaxFailures),
		cooldownDuration: cfg.CooldownDuration,
		lastStateChange:  time.Now(),
	}
	rc.state.Store(int32(StateClosed))
	return rc, nil
}

func (r *RedisCache) isCircuitOpen() bool {
	st := CircuitBreakerState(r.state.Load())
	if st == StateClosed {
		return false
	}
	if st == StateOpen {
		r.mu.Lock()
		defer r.mu.Unlock()
		if time.Since(r.lastStateChange) > r.cooldownDuration {
			r.state.Store(int32(StateHalfOpen))
			r.lastStateChange = time.Now()
			return false
		}
		return true
	}
	return false
}

func (r *RedisCache) recordSuccess() {
	r.consecutiveFailures.Store(0)
	st := CircuitBreakerState(r.state.Load())
	if st != StateClosed {
		r.mu.Lock()
		r.state.Store(int32(StateClosed))
		r.lastStateChange = time.Now()
		r.mu.Unlock()
	}
}

func (r *RedisCache) recordFailure() {
	fails := r.consecutiveFailures.Add(1)
	if fails >= r.maxFailures {
		r.mu.Lock()
		r.state.Store(int32(StateOpen))
		r.lastStateChange = time.Now()
		r.mu.Unlock()
	}
}

func (r *RedisCache) PoolStats() (hits uint64, misses uint64, totalConns uint32, isCircuitOpen bool) {
	stats := r.client.PoolStats()
	return r.hits.Load(), r.misses.Load(), stats.TotalConns, r.isCircuitOpen()
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

	if r.isCircuitOpen() {
		r.misses.Add(1)
		return nil, "", false, nil
	}

	cacheKey := r.generateKey(key, params)
	val, err := r.client.Get(ctx, cacheKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			r.recordSuccess()
			r.misses.Add(1)
			return nil, "", false, nil
		}
		r.recordFailure()
		r.misses.Add(1)
		return nil, "", false, nil
	}

	r.recordSuccess()
	r.hits.Add(1)

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

	if r.isCircuitOpen() {
		return nil
	}

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
		r.recordFailure()
		return nil
	}

	r.recordSuccess()
	return nil
}
