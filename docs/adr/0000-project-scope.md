# ADR-0000: Project Scope

## Status

Accepted

## Context

Optivor is a new open-source project, and without an explicit scope
statement it is at high risk of feature creep. Image infrastructure is a
deep problem space — storage, transformation, caching, deployment,
security, and observability all expand naturally in scope. Left
unconstrained, the first release will never ship, because there is always
"one more provider" or "one more format" to add.

This ADR exists to answer four questions before any code is written:

1. What is Optivor, in one sentence?
2. What is Optivor explicitly *not*?
3. What does the first release (V0) contain, and what does it deliberately
   exclude?
4. How do we know when V0 is done?

Every future scope discussion should be resolved by checking this
document first. If a proposed feature is not in the V0 list below, it is
out of scope for V0 by default, regardless of how small it seems.

## Decision

### What Optivor is

Optivor is an **open-source image infrastructure framework**. It provides
the runtime, driver interfaces, and deployment tooling required to run a
production-grade image pipeline on top of object storage the user already
owns.

### What Optivor is not

- Not an image hosting service (Optivor never stores images itself).
- Not a CDN (Optivor may sit behind one, but does not operate one).
- Not a hosted SaaS product.
- Not a Cloudinary/ImageKit alternative in the sense of being a managed
  API — Optivor is infrastructure the user runs, not a service the user
  calls.

### V0 Scope — Included

V0's only goal is to prove the smallest possible end-to-end pipeline
works and is genuinely portable. V0 includes:

- **Standalone Runtime**: a single Go binary that serves HTTP requests
  and performs image transformations. No CLI required to run it; a
  runtime binary plus a config file is sufficient.
- **Storage**: one storage driver, S3-compatible object storage only
  (validated against both AWS S3 and MinIO to confirm the abstraction
  isn't accidentally AWS-specific).
- **Transformation**: resize only (width/height, with a small set of
  fit modes — cover, contain, fill). No crop, rotate, or watermark in V0.
- **Output formats**: source format passthrough, plus WebP encoding.
  WebP encoding will use a cgo binding to `libwebp` in V0. This is a
  deliberate tradeoff: it sacrifices some cross-compilation simplicity
  in exchange for a mature, correct encoder, rather than depending on an
  immature pure-Go implementation. This tradeoff is recorded, not hidden,
  and should be revisited once pure-Go encoders mature.
- **AVIF is explicitly deferred to V0.1**, not included in V0. As of
  writing, there is no production-quality pure-Go AVIF encoder, and
  adding a second cgo/binary dependency in the first release increases
  build and distribution complexity beyond what V0 needs to prove its
  core thesis. AVIF support should be re-evaluated in ADR form once V0
  ships.
- **Cache**: a basic in-process/filesystem cache only. No Redis, no
  distributed cache, no object-storage-backed cache in V0.
- **Configuration**: a single static config file (format decided in
  ADR-0004 territory, likely YAML), no dynamic policy engine.

### V0 Scope — Explicitly Out of Scope

The following are known future work and are intentionally **not** part
of V0. Listing them here is meant to prevent them from being pulled
forward:

- CLI (`optivor init`, `optivor deploy`, etc.)
- Any deployment adapter (Cloudflare, AWS Lambda, Kubernetes, Fly.io)
- Any storage driver other than S3-compatible (Azure Blob, GCS as
  distinct drivers, even though GCS/Azure may be S3-compatible via
  gateways)
- Plugin system / extension loading of any kind
- Dashboard or any web UI
- AI-based features (smart cropping, content-aware resizing, etc.)
- Watermarking
- Multi-node / distributed runtime, horizontal scaling features
- Signed URL generation and policy-based authorization (tracked
  separately as a security-critical feature for V0.1, not ignored, just
  not in V0)
- Metrics, tracing, structured logging beyond basic stdout logging

### Definition of Done for V0

V0 is complete when a developer can go from an empty S3-compatible
bucket to serving a resized, optionally WebP-converted image through a
single standalone binary, using only a config file, in under five
minutes, with zero cloud-provider-specific configuration and zero
custom code. This statement — not the feature list alone — is the
acceptance test for V0. A feature list can always be argued as
"incomplete"; this statement cannot.

## Consequences

- Any feature request or PR that does not fit inside the V0 list above
  should be labeled `post-v0` and deferred, not silently expanded into
  the current milestone.
- The cgo dependency for WebP is a conscious, documented tradeoff and
  should not be treated as a temporary hack to "fix later" without a
  corresponding ADR when it's revisited.
- Deferring signed URLs and authorization to V0.1 means V0 is **not**
  safe to expose directly to the public internet without an external
  auth layer. This should be stated clearly in the README and V0
  release notes so nobody deploys it unprotected by mistake.
- This scope will likely feel too small to early readers of the
  Architecture Vision document. That is intentional — the Architecture
  Vision describes the long-term shape of the project; this ADR
  describes what ships first.
