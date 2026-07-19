# ADR-0002: High-Level Architecture

## Status

Accepted

## Context

Before any interface or method signature is written, the project needs
agreement on its layer boundaries — what each layer is responsible for,
and just as importantly, what it must never know about. Without this,
early code tends to blur responsibilities (e.g., the runtime accidentally
importing a cloud SDK "just this once"), and those shortcuts become very
expensive to undo later.

This ADR intentionally stays at the level of boxes and arrows. It does
not define Go interfaces, package layouts, or method signatures — those
belong in later, narrower ADRs (ADR-0005 onward) written at the time each
layer is actually implemented.

## Decision

Optivor is organized into four layers, each living in principle in its
own repository, communicating only through narrow, stable boundaries.

```
CLI
  ↓
Runtime
  ↓
Storage Drivers
  ↓
Object Storage (user-owned)
```

and, orthogonally:

```
Deployment Adapter
  ↓
Cloudflare / AWS / Kubernetes / Fly.io / Standalone
```

### CLI layer

Responsible for developer-facing orchestration only: project
initialization, configuration generation, invoking deployment adapters,
and surfacing diagnostics (logs, metrics, health) to the user. The CLI
**never** performs image processing itself and **never** talks to object
storage directly — it delegates to the runtime and to deployment
adapters. This boundary exists so the CLI can evolve (new commands, new
UX) without ever risking the correctness of the image pipeline itself.

### Runtime layer

Responsible for the actual image pipeline: accepting HTTP requests,
routing, applying transformations, and returning responses. The runtime
talks to storage exclusively through the Storage Driver abstraction (see
below) and has **zero knowledge of deployment targets** — it does not
know if it is running on a bare VM, inside a container, or behind a
Cloudflare Worker. This is the layer most contributors will spend time
in, and it is the layer with the strictest "no cloud-provider imports"
rule in the entire codebase.

### Storage Drivers

Responsible for translating a generic "read this object / write this
object" interface into provider-specific API calls (S3 SigV4, etc.). A
storage driver is intentionally narrow in scope — it does not know
anything about image transformation, caching, or HTTP. This narrowness
is what allows a contributor to implement, say, a Backblaze B2 driver by
reading only the driver interface, consistent with the contributor-first
commitment in ADR-0001.

### Deployment Adapters

Responsible for taking the runtime (or, in constrained environments, a
thin proxy in front of the runtime — see the Cloudflare case below) and
making it run on a specific target. Each cloud provider lives in its own
adapter repository. Deployment adapters are explicitly allowed to depend
on provider SDKs and provider-specific tooling (Wrangler, Terraform,
Helm) — this is the one layer where that is acceptable, precisely
because it is isolated from the runtime and storage layers.

**Special case — edge platforms without a full language runtime (e.g.,
Cloudflare Workers):** the Go runtime cannot execute natively inside a
V8-isolate-based platform. For these targets, the deployment adapter
does not run the runtime itself; instead it deploys a thin edge layer
(routing, cache keys, auth, redirects) that forwards actual image
processing to a Go runtime running elsewhere (a VM, container, or
serverless Go target). This distinction — "deploys the runtime" vs.
"deploys a proxy in front of the runtime" — must be made explicit in
each deployment adapter's own documentation, so users are never
surprised by where their images are actually being processed.

## Consequences

- A dependency-direction rule follows directly from this diagram: the
  Runtime must never import a Deployment Adapter or CLI package, and a
  Storage Driver must never import the Runtime. Tooling (e.g., a simple
  import-graph lint check) should eventually enforce this, but the rule
  exists as of this ADR regardless of whether it's automated yet.
- Any deployment target that cannot run the Go runtime natively (edge
  isolates, some serverless platforms) requires its adapter to document
  the "proxy vs. full runtime" distinction explicitly — this is now a
  standing requirement, not a case-by-case judgment call.
- This layering is what makes ADR-0003 (Extensibility Model) possible:
  each layer boundary in this diagram is a candidate extension point.
