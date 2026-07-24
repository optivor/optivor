# ADR-0013: Persistent Caching, Remote Fetching, and Ecosystem Resilience

- **Status:** Accepted
- **Date:** 2026-07-24
- **Authors:** Optivor Team

## Context

Production image optimization services face high costs and performance degradation when deployments invalidate ephemeral local disk caches, when automated web crawlers hit dozens of `srcset` variant URLs simultaneously, or when dynamic remote media fetching lacks domain isolation and SSRF protection. Frontend frameworks like Next.js also require zero-config loaders for native `next/image` integration.

## Decision

1. **Deploy-Proof Persistent Caching (`internal/cache/persistent`)**:
   - Introduce a multi-tier persistent cache driver supporting Redis and S3/Storage backends.
   - Cache keys are generated deterministically using SHA-256 over image source identifier, transformation parameters, and format options.
   - Preserves cached image transformations across application redeployments or container restarts.

2. **Transparent Remote Fetching (`/remote`, `/fetch`)**:
   - Provide `/remote?url=...` and `/fetch?url=...` HTTP endpoints for on-the-fly remote image fetching and processing.
   - Enforce strict domain whitelisting (`allowed_domains` in `optivor.yaml`) and private IP blocking (SSRF guard against 127.0.0.1, 10.0.0.0/8, 169.254.169.254, etc.).

3. **Bot & Crawler Protection / Concurrency Rate Limiting**:
   - Implement User-Agent crawler detection (Googlebot, Bingbot, Bytespider, etc.) and apply dedicated rate limits.
   - Introduce per-variant transformation concurrency semaphores to cap peak CPU spikes when crawlers or responsive design engines request 10-15 image variants simultaneously.

4. **Dynamic Preset Engine (`/preset/:name/*`)**:
   - Support named parameter presets in `optivor.yaml` (e.g., `avatar`, `thumbnail`, `hero`).
   - Requests to `/preset/avatar/photo.jpg` automatically apply configured preset parameters (`w=150,h=150,f=avif,q=80`).

5. **Official Optivor Next.js Image Loader (`@optivor/next`)**:
   - Provide an npm package (`packages/next-loader`) exported as `@optivor/next` for zero-config `next/image` integration.

## Consequences

- Reduced egress and compute costs during redeployments.
- Protection against SSRF vulnerabilities when optimizing external remote images.
- System stability during heavy crawler indexing sweeps.
- Simplified frontend integration for Next.js applications.
