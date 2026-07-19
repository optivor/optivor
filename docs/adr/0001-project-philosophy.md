# ADR-0001: Project Philosophy

## Status

Accepted

## Context

Technical decisions eventually get questioned — "why don't we just add a
hosted mode?", "why not use Node instead of Go?", "why not let people pay
us to store images?". Without a written philosophy, these questions get
re-litigated every time someone new joins the project, and the answer
depends on whoever is in the room that day.

This ADR is not about technology. It is about the principles that every
future technical ADR should be consistent with. If a proposed change
violates this document, the proposal should change, not this document —
unless the philosophy itself is being deliberately revisited, which
should happen rarely and explicitly.

## Decision

### Why open source, and open-source-first (not open-core)

Optivor's value is in being a shared standard, not a proprietary product.
An image pipeline is infrastructure — like Terraform or Caddy — and
infrastructure earns trust through inspectability, not through a vendor's
promises. Optivor will not hold back features for a paid tier. There is
no planned open-core split. If a sustainable business ever emerges around
Optivor, it will be built on top of the open-source project (hosting,
support, managed offerings by third parties), not by withholding
capability from the open-source runtime itself.

### Why "Bring Your Own Storage" (BYOS)

The single biggest source of vendor lock-in in image infrastructure is
storage. If Optivor stored images itself, every user would be one pricing
change away from being trapped. By requiring users to bring their own
S3-compatible bucket, Optivor guarantees that a user's data is never
hostage to Optivor's existence. This is a non-negotiable constraint on
every future feature: **Optivor must never require or default to storing
image bytes it does not already have permission to read from the user's
own bucket.**

### Why not a hosted SaaS

A hosted SaaS product optimizes for retention and lock-in almost by
definition — recurring revenue depends on switching costs. Optivor
optimizes for the opposite: the easier it is to leave, fork, or
self-host, the more trustworthy it is as infrastructure. This is a
deliberate rejection of the Cloudinary/ImageKit business model, not an
oversight.

### Why Go

Go was chosen for three concrete reasons, not just familiarity:

- **Single static binary distribution.** Infrastructure tools that are
  easy to `curl | install` or drop into a Docker `FROM scratch` image
  have historically won developer trust faster than tools requiring a
  language runtime (see: Terraform, Caddy, k6, MinIO).
- **Concurrency model fits the workload.** Image transformation pipelines
  are I/O-bound (storage fetch) interleaved with CPU-bound work
  (encoding), and goroutines/channels map cleanly onto that without the
  callback or async/await complexity of Node or Python equivalents.
- **Ecosystem precedent.** The infrastructure tools Optivor wants to be
  compared against — Terraform, Caddy, Vault, MinIO, k6 — are
  overwhelmingly written in Go. This is not cargo-culting; it means
  contributors coming from that ecosystem already have the right mental
  model and tooling.

### Why cloud/provider-agnostic by default

Every dependency on a specific cloud provider is a dependency Optivor
cannot control the pricing, availability, or roadmap of. The runtime
must always be deployable on a plain VM with nothing but a binary and a
config file — every deployment adapter (Cloudflare, AWS, Kubernetes) is
additive convenience on top of that baseline, never a requirement for it
to function.

### Why contributor-first, and what that actually means

"Contributor-first" is not a slogan; it means a specific architectural
commitment: **a contributor should be able to own an entire subsystem
without needing to understand the rest of the codebase.** Someone who
wants to add Azure Blob support should be able to read the Storage
Driver interface, implement it, write tests against it, and submit a PR
— without reading the CLI code, the deployment adapters, or the cache
layer. This is why Optivor is deliberately organized into layers with
narrow, stable interfaces between them (see ADR-0002 and ADR-0003) rather
than a single monolithic codebase where every change touches everything.
This is also why the project favors composition and interfaces over deep
inheritance or shared mutable state — coupling is the direct enemy of
this goal.

### Why not feature-richness as the primary goal

Optivor will consistently lose feature-for-feature comparisons against
managed products like Cloudinary, which can move faster because they
don't need cross-provider portability, community review, or backward
compatibility guarantees. Optivor's goal is not to win that comparison.
Its goal is to be the simplest, most portable, most extensible foundation
available — a foundation other tools and companies build on top of,
rather than a finished product competing feature-for-feature with SaaS
alternatives.

## Consequences

- Any proposal that introduces hosted storage, a paid feature tier, or a
  hard dependency on a single cloud provider should be treated as a
  philosophy-level change requiring explicit, visible discussion — not a
  normal PR.
- Contributors should be able to point to this document when justifying
  why a PR was structured a certain way (e.g., "the Storage Driver
  interface stays minimal because ADR-0001 commits us to
  subsystem-level ownership").
- This document should be revisited if the project's fundamental
  situation changes (e.g., a foundation adopts the project, or a
  commercial entity forms around it) — but any such revision should be
  its own ADR, not a silent edit to this one.
