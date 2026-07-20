# ADR-0006: V0.2 Deployment Adapter Selection

## Status

Accepted

## Context

Per **ADR-0002** (High-Level Architecture) and **ADR-0003** (Extensibility Model), Deployment Adapters exist to package or deploy the Optivor runtime onto target platforms.
ADR-0002 establishes a strict distinction between:
1. **Full-Runtime Adapters**: Deploy the complete Optivor Go binary (and C-libraries like `libvips`) directly to target infrastructure (e.g. systemd, VMs, containers).
2. **Proxy Adapters**: Deploy an edge proxy (e.g. Cloudflare Worker or AWS Lambda) that routes requests to an upstream Optivor runtime binary.

For the **V0.2 milestone**, Optivor requires its first real Deployment Adapter to enable reproducible, production-grade deployments without forcing users to write custom init scripts or process supervisors.

## Options Considered

1. **Systemd Service Packaging (`systemd`)**:
   Package the Optivor binary with systemd unit file templates, `nfpm` `.deb`/`.rpm` packages, and Makefile targets.
   - *Pros*: Zero external cloud platform dependency, highly reliable, standard Linux deployment model, completely self-contained.
   - *Cons*: Limited to Linux OS distributions supporting systemd.

2. **Fly.io Adapter (`flyio`)**:
   Provide a Fly.io Docker container deployment adapter via `fly.toml` scaffolding and Fly CLI invocation.
   - *Pros*: Easy cloud hosting with persistent disk cache support.
   - *Cons*: Introduces dependency on third-party cloud platform APIs during early release cycles.

3. **Cloudflare Worker Edge Proxy Adapter (`cloudflare`)**:
   Deploy a Cloudflare Worker proxy in front of an upstream runtime.
   - *Pros*: Excellent global caching and edge auth.
   - *Cons*: High architectural complexity due to the proxy-vs-runtime dual deployment model.

## Decision

We decide to adopt **Systemd Service Packaging** as the initial Deployment Adapter for V0.2.

Key architectural boundaries and implementation choices:
- Scaffolding unit template file at `deploy/systemd/optivor.service`.
- Packaging support for Debian/Ubuntu (`.deb`) and RHEL/Fedora (`.rpm`) via GoReleaser (`nfpm`).
- `Makefile` installation helper targets (`make install`, `make uninstall`).
- CLI orchestration of the systemd deployment adapter via out-of-process execution (ADR-0003 boundary).

Fly.io and Cloudflare deployment adapters are deferred to V0.3+ / V1 milestones.

## Consequences

- Optivor maintains its core philosophy: **Bring Your Own Infrastructure, zero cloud vendor lock-in**.
- Linux server deployments can run Optivor as a managed background daemon with automatic restart on failure (`Restart=always`).
- System logs integrate natively with `journalctl`, providing structured logging capability.
