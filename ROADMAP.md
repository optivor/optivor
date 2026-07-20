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

## V0.6 — Multi-Provider & Multi-Bucket Management (Planned)

Goal: Manage multiple storage buckets across diverse providers (AWS S3, Cloudflare R2, Backblaze B2, GCS) from a single unified Optivor instance.

- [ ] Multi-bucket routing engine (`internal/storage/router`) mapping bucket identifiers/prefixes to specific driver instances
- [ ] Declarative multi-bucket configuration schema in `optivor.yaml`
- [ ] Cross-provider failover and primary/backup bucket fallback policies
- [ ] Per-bucket metrics and telemetry breakdown in `optivor metrics` and OpenTelemetry spans
- [ ] CLI multi-bucket verification in `optivor doctor`

## V1 — Extension points

- Storage Driver interface finalized and documented for external
  contributors (additional drivers beyond S3-compatible)
- Runtime Module mechanism decided (see `docs/adr/0003-extensibility-model.md`)
  and documented
- Additional Deployment Adapters (Cloudflare, AWS, Kubernetes) — each
  documenting explicitly whether it deploys the full runtime or a proxy
  in front of it, per ADR-0002

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
