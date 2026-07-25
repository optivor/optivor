# Roadmap

This roadmap is intentionally coarse. Detailed scope for the current
milestone lives in [`docs/adr/0000-project-scope.md`](./docs/adr/0000-project-scope.md)
— this file exists to show the overall trajectory and to prevent
premature work on later milestones.

## V0 — Prove the core loop (Completed)

Goal: a single standalone binary can serve a resized, WebP-converted
image from an S3-compatible bucket, with no cloud-specific configuration.

Full scope: [`docs/adr/0000-project-scope.md`](./docs/adr/0000-project-scope.md).

Status: ✅ Completed.

## V0.1 — Make it safe to expose (Completed)

- [x] Signed URLs and basic policy-based authorization (ADR-0005)
- [x] AVIF output support (`libheif-dev`)
- [x] Resource limits: max pixel count, max memory per transform, request timeouts (decompression-bomb protection)
- [x] Basic per-IP rate limiting
- [x] LRU cache eviction
- [x] Cross-compilation matrix (Zig-based toolchain in `.goreleaser.yaml`)

## V0.2 — First deployment adapter (Completed)

- [x] Standalone-as-systemd-service packaging deployment adapter (ADR-0006)
- [x] Structured logging (JSON/Text slog) and basic metrics (Prometheus endpoint)

## V0.3 — CLI (Completed)

- [x] `optivor init`, `optivor deploy`, config scaffolding
- [x] CLI orchestrates systemd deployment adapter (ADR-0007)

## V0.4 — Observability (Completed)

- [x] OpenTelemetry tracing (ADR-0008, stdout + OTLP gRPC)
- [x] `optivor doctor`, `optivor logs`, `optivor metrics` diagnostics commands

## V0.5 — Docker-First & Provider Agnostic Architecture (Completed)

- [x] Docker-first deployment strategy (ADR-0009) with volume-mounted config and `--provider` flag
- [x] Provider Registry & driver binary convention (ADR-0010)
- [x] `optivor driver install` CLI subcommand suite for external driver management
- [x] Complete documentation wiki (`docs/wiki/`) and developer driver creation guides

## V0.6 — Multi-Provider & Multi-Bucket Management (Completed)

Goal: Manage multiple storage buckets across diverse providers (AWS S3, Cloudflare R2, Backblaze B2, GCS) from a single unified Optivor instance.

- [x] Multi-bucket routing engine (`internal/storage/router`) mapping bucket identifiers/prefixes to specific driver instances (ADR-0011)
- [x] Declarative multi-bucket configuration schema in `optivor.yaml`
- [x] Cross-provider failover and primary/backup bucket fallback policies
- [x] Per-bucket metrics and telemetry breakdown in `optivor metrics` and OpenTelemetry spans
- [x] CLI multi-bucket verification in `optivor doctor`

## V0.6.1 — GitHub Driver Installation (Completed)

- [x] `optivor driver install` support for `github:org/repo[@tag]` shorthand and HTTPS URLs

## V0.7 — Edge Integration & Lifecycle Management (Completed)

- [x] Edge deployment strategy and proxy routing guide (ADR-0012)
- [x] CLI `optivor bucket lifecycle` (`list`, `set`, `delete`) management subcommands

## V0.8 — Production Hardening & Benchmark Gating (Completed)

- [x] External driver protocol SDK specification (`docs/wiki/driver-sdk-specification.md`)
- [x] CI benchmark regression gating job (`go test -bench`)

## V0.9 — Production Resilience & Developer Ecosystem (Completed)

Goal: Address real-world Next.js & Crawler cost spikes, deployment cache invalidation issues, and provide zero-config frontend ecosystem integrations.

- [x] Deploy-Proof Persistent Caching (Redis/Storage-backed cache layer preserved across application redeployments)
- [x] Transparent Remote Fetching (`/remote` & `/fetch` endpoints with strict domain whitelist for on-the-fly optimization)
- [x] Bot & Crawler Protection / Smart Throttling (Rate-limiting and variant concurrency guards against automated web crawlers hitting 10-15 `srcset` URLs)
- [x] Dynamic Preset Engine (Named URL alias presets e.g., `/preset/avatar/photo.jpg` -> `w=150,h=150,f=avif,q=80`)
- [x] Official Optivor Next.js Image Loader (`@optivor/next` npm package for zero-config `next/image` integration)


## V1 — Production Readiness & Extension Points (Completed)

Goal: Professional-grade horizontal scaling, B2B deployment primitives, smart media optimization, and stable extensibility APIs.

Status: ✅ Completed.

- [x] **Smart Cropping (`fit=smart`)**: Entropy/attention-based cropping utilizing libvips `InterestingAttention`.
- [x] **Redis Cache Backend**: Stateless and horizontally scalable caching layer across multiple pods/nodes.
- [x] **Helm Chart & Kubernetes manifests**: Production-ready helm installation and configuration values for B2B/Cloud deployments.
- [x] Storage Driver interface finalized and documented for external contributors (additional drivers beyond S3-compatible)
- [x] Runtime Module mechanism decided (ADR-0015) and documented
- [x] Additional Deployment Adapters (Cloudflare Worker Edge Proxy, AWS ECS Fargate Container) — each documenting explicitly whether it deploys the full runtime or a proxy in front of it, per ADR-0002
- [x] **Enterprise IAM & Dynamic Access Control**:
  - [x] Support AWS IRSA (IAM Roles for Service Accounts) and GKE Workload Identity natively (token credential providers instead of static keys).
  - [x] Implement client-level API-Key token authorization policies with bucket-level operation scopes (read/write/lifecycle).
  - [x] Implement Dynamic Signed URL delegation with HMAC secret validation for private bucket isolation.
- [x] **Production Infrastructure Hardening**:
  - [x] Add `NetworkPolicy` to the Helm Chart to restrict outbound pod traffic to whitelisted endpoints.
  - [x] Implement Redis connection pool sizing controls and a circuit-breaker for graceful S3 fallback when Redis goes offline.
  - [x] Add custom Prometheus `/metrics` endpoint (exposing request latencies, libvips cache utilization, and Redis pool stats).
  - [x] Add `checksum/config` annotation to the Helm Deployment pod template for automatic rolling config updates.

## V1.1 — Developer Experience & Zero-Friction Onboarding (Completed)

Goal: Eliminate developer setup friction, enable zero-config local trials, and provide 1-line installation tools.

Status: ✅ Completed.

- [x] **Zero-Config Local Storage Fallback Mode**: Allow `optivor start` to serve local disk images (`./storage`) out-of-the-box when no YAML or S3 credentials are provided.
- [x] **100% Environment Variable Configuration**: Support zero-file configuration (`OPTIVOR_S3_BUCKET`, `OPTIVOR_S3_REGION`, etc.) for seamless PaaS/Serverless container execution.
- [x] **1-Line Shell Installer**: `curl -fsSL https://optivor.app/install.sh | bash` script auto-detecting OS/architecture and adding the Optivor binary to system `PATH`.
- [x] **Interactive CLI Wizard (`optivor init --interactive`)**: Step-by-step CLI prompt to configure S3/R2 credentials, test connection live, and generate `optivor.yaml`.
- [x] **1-Click PaaS Deployment Templates**: Deploy buttons and templates for Railway, Render, and Fly.io in the primary documentation and README.
- [x] **ADR & Architecture Specifications**:
  - [x] `docs/adr/0016-zero-config-fallback-and-env-override-matrix.md` — ADR specifying zero-config local fallback mechanics and environment variable override precedence.
- [x] **Deployment Guides & Wiki Documentation**:
  - [x] `docs/deployment/paas-railway-render-fly.md` — 1-Click deployment guides and container blueprints for Railway, Render, and Fly.io.
  - [x] `docs/wiki/zero-config-quickstart.md` — Guide for running Optivor locally in under 10 seconds without cloud credentials.
  - [x] `docs/wiki/cli-wizard-guide.md` — Detailed CLI reference and walkthrough for `optivor init --interactive`.

## V1.2 — Advanced Engine Capabilities & Ecosystem SDKs (Planned)

Goal: Expand dynamic media transformation features and extend frontend/backend framework integrations.

- [ ] **Dynamic Watermarking & Overlays**: Image overlays with custom position (`gravity`), opacity, and scaling controls (`overlay=logo.png&gravity=bottom_right&opacity=50`).
- [ ] **Manual Focal Point Cropping**: Direct coordinate-based focal cropping (`focal=0.3,0.7`) to complement smart entropy cropping.
- [ ] **Animated Media Conversion**: GIF to animated WebP/MP4 micro-conversions for bandwidth optimization.
- [ ] **Image Filter Effects**: Dynamic blur (`blur=10`), grayscale (`grayscale=true`), and pixelation (`pixelate=5`).
- [ ] **Ecosystem Client SDKs**: Official packages for `@optivor/react`, `@optivor/vue` (Nuxt), `@optivor/js`, `optivor-php` (Laravel), and `optivor-python` (Django/FastAPI).
- [ ] **ADR & Architecture Specifications**:
  - [ ] `docs/adr/0017-watermarking-overlays-and-focal-crop.md` — Architectural specification for libvips overlay compositing and focal point calculations.
- [ ] **Transformation Guides & SDK Specifications**:
  - [ ] `docs/wiki/watermarking-and-effects.md` — Comprehensive parameter reference for watermarks, overlays, and image filters.
  - [ ] `docs/wiki/client-sdk-specification.md` — Standardized framework SDK contract for official client libraries.

## Future Planned — Granular IAM & Path-Level Access Control (Target Version TBD)

Goal: Provide fine-grained role-based access control (RBAC) and IAM path policies for multi-bucket / single URL routing, enabling user-level path isolation (e.g., user-A restricted to `s3://bucket/folder-a/*`).

- [ ] **Path-Based Prefix Authorization Policies**: Restrict API key tokens and signed URLs to specific object key prefixes/folders within S3/R2 buckets.
- [ ] **Role-Based Privilege Definition (IAM Rules)**: Support declarative role definitions (`admin`, `editor`, `reader-path-only`) mapping identities to bucket/prefix capabilities.
- [ ] **ADR & Architecture Specifications**:
  - [ ] `docs/adr/0019-iam-path-privilege-and-role-policy.md` — Specification for IAM role evaluation and path-level prefix validation.
- [ ] **Wiki Documentation**:
  - [ ] `docs/wiki/iam-path-authorization-guide.md` — Guide for defining folder-level privileges and role policies in Optivor.

## Future Planned — Open-Core Cloud Control Plane Integration (Target Version TBD)

Goal: Allow self-hosted Optivor Engine instances to optionally pair with the free `optivor.app` SaaS Cloud Dashboard for real-time analytics, bandwidth savings tracking, and visual monitoring.

- [ ] **Opt-in Cloud Pairing Engine (`optivor connect`)**: Add optional non-blocking background telemetry sync daemon connecting self-hosted instances via `OPTIVOR_CLOUD_KEY` / Access Key pairing.
- [ ] **Privacy-First Aggregate Telemetry Protocol**: Transmit non-sensitive aggregate operational metrics (request count, cache HIT/MISS ratio, bandwidth saved, latency distribution) — zero image binaries or private payloads transmitted.
- [ ] **Fail-Safe Telemetry Fallback**: Guarantee 100% uninterrupted local image processing when cloud control plane network connection is unavailable or unconfigured.
- [ ] **ADR & Architecture Specifications**:
  - [ ] `docs/adr/0018-cloud-control-plane-pairing-and-telemetry-protocol.md` — ADR specifying non-blocking phone-home telemetry protocol, fail-safe isolation, and privacy guarantees.
- [ ] **Wiki Documentation**:
  - [ ] `docs/wiki/cloud-dashboard-pairing.md` — Walkthrough for pairing self-hosted Optivor instances with the free `optivor.app` dashboard.

## Explicitly not scheduled

These are known ideas that are **not** on this roadmap yet. They are not
rejected — they simply haven't been through the scope discipline
described in `CONTRIBUTING.md`:

- Dashboard / web UI
- AI-based transformations (smart cropping, content-aware resizing)
- Multi-node / horizontally scaled runtime
- Formal governance structure (RFC process, TSC) — revisit once the
  contributor base and organizational backing justify it (see the
  governance note in `docs/adr/0001-project-philosophy.md`)

If you want to propose moving something from this list onto an active
milestone, open a Feature Request issue — see `CONTRIBUTING.md`.
