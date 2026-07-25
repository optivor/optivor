# ADR-0015: Runtime Module Mechanism & Extension Architecture

- **Status**: Accepted
- **Deciders**: Optivor Core Architecture Team
- **Date**: 2026-07-25

---

## Context and Problem Statement

Optivor's extensibility model (ADR-0003) distinguished between out-of-process Storage Drivers, out-of-process Deployment Adapters, and in-process **Runtime Modules** (transforms, cache backends, auth/policy handlers).

As Optivor reaches V1 production readiness, we require a formal decision on how Runtime Modules are structured, compiled, registered, and executed within the runtime without introducing cgo/plugin instability or dynamic shared object overhead.

---

## Decision Outcome

We decide to adopt **In-Process Interface Contracts & Declarative Compile-Time Registry** for all Runtime Modules:

1. **Transform Modules**: Implement the `pipeline.Transformer` interface interface contracts. Custom transform handlers register with `pipeline.RegisterTransform(name, handler)`.
2. **Cache Backend Modules**: Implement the `cache.Cache` interface contract (`Get`, `Set`, `Close`, optional `PoolStats`). Custom backends register with `cache.RegisterBackend(name, factory)`.
3. **Auth & Access Control Modules**: Implement the `server.AuthPolicy` interface contract (`ValidateRequest(r *http.Request) error`).

Dynamic C-Go `.so` shared library plugins are explicitly **rejected** due to memory safety risks, libvips thread synchronization complexities, and lack of cross-platform reproducibility. Instead, Optivor uses Go interface composition with static build-tag compilation for custom runtime extensions.

---

## Architecture & Interfaces

### 1. Cache Module Interface
```go
type Cache interface {
    Get(ctx context.Context, key string, params TransformParams) ([]byte, string, bool, error)
    Set(ctx context.Context, key string, params TransformParams, data []byte, contentType string) error
    Close() error
}
```

### 2. Transform Module Interface
```go
type TransformModule interface {
    Name() string
    Apply(img *vips.ImageRef, params TransformParams) error
}
```

### 3. Auth & Policy Module Interface
```go
type PolicyModule interface {
    Name() string
    Authorize(r *http.Request, bucket string, scope string) bool
}
```

---

## Consequences

### Positive
- Zero runtime overhead or RPC latency for image transformations and cache lookups.
- Type-safe, thread-safe, and memory-safe execution inside the main binary.
- Clean layer separation matching ADR-0002.

### Negative
- Custom transforms or custom cache backends require re-compiling the Optivor binary (standard Go module extension pattern).
