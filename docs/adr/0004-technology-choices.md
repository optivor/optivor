# ADR-0004: Technology Choices

## Status

Accepted

## Context

Early on, someone will inevitably suggest an alternative to a core
technology choice — "why not Node for the CLI," "why not Fiber instead
of chi," "why not Zap instead of slog." Most of these are individually
reasonable suggestions, but re-litigating them one at a time wastes
review time and produces an inconsistent stack. This ADR records the
baseline technology choices and enough rationale that a future proposal
to change one of them can be evaluated against the reason it was chosen,
rather than argued from scratch.

This is intentionally the lightest-weight ADR in the initial set — these
are reversible, low-blast-radius decisions compared to ADR-0000 through
ADR-0003. It exists mainly to save review time, not to enshrine
permanence.

## Decision

| Concern | Choice | Rationale |
|---|---|---|
| Core language | Go | Static binary distribution, concurrency model fits I/O+CPU-bound image pipeline workload, dominant language in comparable infrastructure tools (Terraform, Caddy, MinIO, Vault) — see ADR-0001. |
| CLI framework | Cobra | De facto standard for Go CLIs; used by Kubernetes, Docker, Hugo, GitHub CLI. Minimizes the learning curve for contributors coming from that ecosystem. |
| Configuration | Viper | Pairs naturally with Cobra; supports layered config (file, env, flags) without custom parsing code. |
| HTTP layer | `net/http` + chi | `chi` is a thin, idiomatic router built directly on `net/http`, not a full framework — consistent with the "minimal dependencies" principle in the Architecture Vision. Avoids pulling in a heavier framework (Fiber, Gin) whose non-standard `Context` types would leak framework-specific assumptions into the runtime. |
| Logging | `log/slog` (standard library) | Structured logging without adding a third-party dependency; sufficient for V0's basic logging needs (see ADR-0000) and forward-compatible with more advanced observability work later. |
| Metrics | Prometheus client library | Industry-standard for infrastructure tooling; pull-based model fits a self-hosted runtime better than push-based alternatives for V0's use case. |
| Tracing | OpenTelemetry | Vendor-neutral tracing standard, consistent with the project's anti-lock-in philosophy (ADR-0001) — this choice is really a philosophy-driven one disguised as a technology one. |
| Testing | Standard library `testing` package | No need for a third-party assertion/mocking framework at this stage; keeps the dependency graph minimal for contributors. Revisit only if table-driven tests become unwieldy. |
| Build & release | GoReleaser | Handles cross-compilation and multi-channel distribution (Homebrew, APT, RPM, Docker, GitHub Releases) with a single config file, directly supporting the "single static binary, easy install" goal in ADR-0001. |

## Consequences

- A proposal to replace any row in this table should reference the
  specific rationale it's arguing against (e.g., "chi's minimalism is no
  longer sufficient because X"), not just present the alternative in
  isolation.
- Because these are lower-stakes, reversible decisions, this ADR can be
  amended more freely than ADR-0000 through ADR-0003 — a superseding PR
  updating this table is sufficient; it does not require a new
  standalone ADR number unless the change is contentious enough to
  warrant its own discussion record.
- Notably absent from this table: the Storage Driver interface
  technology, the Deployment Adapter IPC mechanism, and the Runtime
  Module extension mechanism. Those are deliberately excluded here and
  left to ADR-0005 and later, per ADR-0003.
