package router

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/optivor/optivor/internal/storage"
)

var (
	ErrBucketNotFound = errors.New("bucket alias not found")
	ErrCircularFallback = errors.New("circular fallback chain detected")
)

type BucketTarget struct {
	Alias    string
	Driver   storage.StorageDriver
	Policy   AccessPolicy
	Fallback string
	Provider string
}

type BucketRouter interface {
	Resolve(ctx context.Context, alias string) (storage.StorageDriver, string, error)
	Policy(alias string) AccessPolicy
	Provider(alias string) string
	Aliases() []string
	Target(alias string) (BucketTarget, bool)
}

type DefaultRouter struct {
	targets map[string]BucketTarget
	defaultAlias string
	mu      sync.RWMutex
}

func NewDefaultRouter(targets map[string]BucketTarget, defaultAlias string) (*DefaultRouter, error) {
	r := &DefaultRouter{
		targets: targets,
		defaultAlias: defaultAlias,
	}

	if err := r.validate(); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *DefaultRouter) validate() error {
	for alias, target := range r.targets {
		visited := map[string]bool{alias: true}
		curr := target.Fallback
		for curr != "" {
			if visited[curr] {
				return fmt.Errorf("%w: %s -> %s", ErrCircularFallback, alias, curr)
			}
			visited[curr] = true
			t, exists := r.targets[curr]
			if !exists {
				return fmt.Errorf("fallback alias %q for %q does not exist", curr, alias)
			}
			curr = t.Fallback
		}
	}
	return nil
}

func (r *DefaultRouter) Target(alias string) (BucketTarget, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if alias == "" {
		alias = r.defaultAlias
	}
	t, ok := r.targets[alias]
	return t, ok
}

func (r *DefaultRouter) Policy(alias string) AccessPolicy {
	t, ok := r.Target(alias)
	if !ok {
		return PolicyPublic
	}
	return t.Policy
}

func (r *DefaultRouter) Provider(alias string) string {
	t, ok := r.Target(alias)
	if !ok {
		return "s3"
	}
	if t.Provider != "" {
		return t.Provider
	}
	return "s3"
}

func (r *DefaultRouter) Aliases() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	aliases := make([]string, 0, len(r.targets))
	for a := range r.targets {
		aliases = append(aliases, a)
	}
	return aliases
}

func (r *DefaultRouter) Resolve(ctx context.Context, alias string) (storage.StorageDriver, string, error) {
	if alias == "" {
		alias = r.defaultAlias
	}
	r.mu.RLock()
	target, exists := r.targets[alias]
	r.mu.RUnlock()

	if !exists {
		return nil, "", ErrBucketNotFound
	}

	return &FailoverDriver{
		primaryAlias: alias,
		primaryDriver: target.Driver,
		router: r,
	}, alias, nil
}

type FailoverDriver struct {
	primaryAlias  string
	primaryDriver storage.StorageDriver
	router        *DefaultRouter
}

func (fd *FailoverDriver) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	currAlias := fd.primaryAlias
	driver := fd.primaryDriver

	for currAlias != "" && driver != nil {
		rc, err := driver.Get(ctx, key)
		if err == nil {
			return rc, nil
		}

		target, exists := fd.router.Target(currAlias)
		if !exists || target.Fallback == "" {
			return nil, err
		}

		currAlias = target.Fallback
		nextTarget, ok := fd.router.Target(currAlias)
		if !ok {
			return nil, err
		}
		driver = nextTarget.Driver
	}

	return nil, storage.ErrNotFound
}
