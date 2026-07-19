# Roadmap

This roadmap is intentionally coarse. Detailed scope for the current
milestone lives in [`docs/adr/0000-project-scope.md`](./docs/adr/0000-project-scope.md)
— this file exists to show the overall trajectory and to prevent
premature work on later milestones.

## V0 — Prove the core loop (current)

Goal: a single standalone binary can serve a resized, WebP-converted
image from an S3-compatible bucket, with no cloud-specific configuration.

Full scope: [`docs/adr/0000-project-scope.md`](./docs/adr/0000-project-scope.md).

Status: 🚧 in progress — architecture decisions locked (ADR-0000 through
ADR-0004), runtime implementation not yet started.

## V0.1 — Make it safe to expose

- Signed URLs and basic policy-based authorization
- AVIF output (pending a decision on encoder dependency — see the note
  in ADR-0000)
- Resource limits: max pixel count, max memory per transform, request
  timeouts (decompression-bomb protection)
- Basic rate limiting

## V0.2 — First deployment adapter

- One real Deployment Adapter (target to be decided — candidates:
  standalone-as-systemd-service packaging, or Fly.io, chosen for
  simplicity relative to Cloudflare's proxy-vs-runtime complexity — see
  ADR-0002)
- Structured logging and basic metrics (Prometheus)

## V0.3 — CLI

- `optivor init`, `optivor deploy`, config scaffolding
- CLI orchestrates the Deployment Adapter(s) that exist by this point

## V0.4 — Observability

- OpenTelemetry tracing
- `optivor doctor`, `optivor logs`, `optivor metrics` diagnostics commands

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
