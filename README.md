# Optivor

> Every engineering team eventually builds an image pipeline. Optivor exists so they don't have to build it twice.

## What is Optivor?

Optivor is an **open-source image infrastructure framework**.

It is not an image hosting service. It does not store your images. It is
not a CDN, and it does not operate one. It is not a hosted product you
sign up for.

Optivor is the runtime, driver interfaces, and deployment tooling that
let you run a production-grade image pipeline on top of object storage
you already own.

**Bring Your Own Storage. Control Everything Else.**

## Why Optivor Exists

Almost every product with users ends up needing an image pipeline —
product photos, avatars, chat attachments, generated assets. And almost
every team builds roughly the same thing: fetch from storage, resize,
convert to a modern format, cache, serve, repeat. Then they rebuild it
again at the next company, or the next project.

This isn't a hard problem because it's conceptually novel. It's a hard
problem because it's tedious, easy to get subtly wrong (cache
invalidation, format negotiation, resource limits), and rarely worth a
team's full attention — so it gets built once, badly, and then lived
with for years.

Optivor exists to be the version of this infrastructure you don't have
to build yourself, without asking you to give up ownership of your data
to get it.

## Core Principles

- **Bring Your Own Storage.** Optivor never stores your images. Your
  data lives in your bucket, under your account, under your control —
  always.
- **Open-source-first, not open-core.** There is no paid tier hiding
  behind the free one. What you see in this repository is the product.
- **Provider-agnostic by default.** The runtime works on a plain VM with
  nothing but a binary and a config file. Cloud-specific deployment is
  optional convenience layered on top, never a requirement.
- **Composable, not monolithic.** Storage, transformation, caching, and
  deployment are independently replaceable. You should be able to swap
  any one of them without touching the others.
- **Contributor-first.** You should be able to own an entire piece of
  this project — a storage driver, a deployment adapter — without
  needing to understand the whole codebase to do it.

## Architecture Overview

```
CLI
  ↓
Runtime
  ↓
Storage Drivers
  ↓
Object Storage (yours)
```

```
Deployment Adapter
  ↓
Cloudflare / AWS / Kubernetes / Fly.io / Standalone
```

The runtime knows nothing about cloud providers. Deployment adapters
know nothing about image processing. Storage drivers know nothing about
either. Each piece does one job.

The full reasoning behind these boundaries — and the architectural
decisions that shape them — lives in [`docs/adr/`](./docs/adr), starting
with [ADR-0000: Project Scope](./docs/adr/0000-project-scope.md).

## Getting Started

> Optivor is in early development. The commands below describe the goal
> for the first release (V0), not something you can run today. Follow
> [`ROADMAP.md`](./ROADMAP.md) for real status.

```bash
# point Optivor at an S3-compatible bucket you already own
optivor-runtime --config optivor.yaml
```

No account. No dashboard. No dependency on Optivor's continued
existence. Just a binary, your bucket, and a config file.

## Roadmap

See [`ROADMAP.md`](./ROADMAP.md) for the current milestone and what's
explicitly out of scope for now.

## Contributing

Optivor is designed so a contributor can own an entire subsystem — a
storage driver, a deployment adapter, a runtime module — without needing
to understand the rest of the codebase. See
[`CONTRIBUTING.md`](./CONTRIBUTING.md) to get started.

## License

Apache License 2.0. See [`LICENSE`](./LICENSE).
