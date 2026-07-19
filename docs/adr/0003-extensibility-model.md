# ADR-0003: Extensibility Model

## Status

Accepted

## Context

The Architecture Vision document uses the word "plugin" for almost
everything extensible — storage, deployment, transforms, cache, auth.
Treating all of these as one undifferentiated "plugin system" is a
mistake: they have fundamentally different lifecycles, trust
requirements, and performance constraints, and a single implementation
mechanism cannot serve all of them well.

This ADR exists to name the extension categories precisely and to
record architectural *decisions* about their constraints — without
prematurely locking in specific technologies (gRPC, WASM, etc.) before
the runtime exists to validate those choices against. Technology
selection for each category belongs in its own later ADR (e.g., ADR-0005
for Storage Driver API, ADR-0009 for Runtime Module mechanism), written
once there is real code to validate assumptions against.

## Decision

Optivor recognizes three distinct extension categories. Using the wrong
name for any of these in code, docs, or discussion should be treated as
a bug to fix, because the terminology is load-bearing — it's how
contributors will know which rules apply to what they're building.

### 1. Storage Driver

Implements read/write access to a specific object storage backend (S3,
Azure Blob, GCS, Backblaze B2, MinIO, etc.). Storage Drivers are
data-path components: they are invoked on every request that misses
cache, so they must be low-latency and are expected to run **in-process**
with the runtime.

**Decision:** Storage Drivers are compiled into the runtime binary (or a
build variant of it) rather than loaded dynamically at runtime. This
avoids introducing IPC/RPC latency into the hot data path and avoids the
well-known fragility of Go's native `plugin` package (version/build-flag
coupling with the host binary). Community-contributed drivers are
expected to be built via a build-tag or driver-registration mechanism
decided in a future, narrower ADR — not via dynamic loading.

### 2. Deployment Adapter

Takes the runtime and makes it run on a specific target (Cloudflare, AWS,
Kubernetes, Fly.io, standalone). Deployment Adapters run at deploy time
only, not at request time, and are inherently provider-specific — they
are expected to depend on heavyweight, provider-specific SDKs and CLIs
that must never leak into the runtime or storage layers.

**Decision:** Deployment extensions must execute **out-of-process** from
the CLI, to isolate provider-specific dependencies (SDKs, credentials,
CLI tool versions) from the core CLI binary and from each other. This is
a decision about isolation, not about a specific IPC technology — the
mechanism (gRPC, plain subprocess + stdin/stdout JSON, or something else)
is left to a future ADR informed by real deployment-adapter code. The
precedent this follows (Terraform, Vault) is cited for validation, not
copied uncritically.

### 3. Runtime Module

Extends the request-time behavior of the runtime itself: transformations
(resize, crop, watermark, format encoders), cache backends, and
authentication/policy checks. Runtime Modules run in the hot path, on
every request or every cache miss, so their performance and safety
characteristics matter far more than Deployment Adapters'.

**Decision:** The implementation mechanism for Runtime Modules is
explicitly **not decided by this ADR**. Candidates include in-process Go
interfaces with build-time registration (fastest, least isolated),
sandboxed WASM modules (safer, some performance cost, enables
runtime-loadable community plugins without recompiling), or an
RPC-per-request model (rejected as a *default* due to expected latency
cost for image workloads, though it may be acceptable for non-hot-path
Runtime Modules like custom auth checks). This decision is deferred to a
dedicated ADR written once the core transform pipeline exists and can be
benchmarked against real workloads — but the category and its
constraints (must run in or very near the hot path; must not be allowed
to compromise the isolation the runtime otherwise guarantees) are fixed
as of this ADR.

### The general principle

> Not every extension mechanism uses the same implementation technology.

Extension category is determined by **where in the request/deploy
lifecycle the extension runs**, not by convenience of implementation.
Committing to a single mechanism (e.g., "everything is a gRPC plugin," or
"everything is WASM") for all three categories would either cripple
hot-path performance (Storage Drivers, Runtime Modules) or force
unnecessary isolation overhead onto deploy-time-only code (Deployment
Adapters). Each category's mechanism should be chosen to fit its
lifecycle, not the other way around.

## Consequences

- Documentation, issue templates, and CONTRIBUTING.md should always use
  these three terms precisely — "plugin" as a generic word should be
  avoided in technical contexts from this point forward.
- ADR-0005 (Storage Driver API) and a future ADR for Runtime Module
  mechanism become required reading before any external contributor
  starts a new driver or transform — this ADR is the map, not the
  territory.
- Because the Runtime Module mechanism is deliberately left open, any
  early transform code (e.g., the V0 resize/WebP path from ADR-0000)
  should be written as a plain internal Go function first, not
  prematurely wrapped in an extension abstraction that hasn't been
  decided yet. Over-engineering the extension point before it's needed
  would violate the "decisions over solutions" spirit of this ADR.
