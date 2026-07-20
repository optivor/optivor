# Introduction to Optivor

Optivor is an **open-source image infrastructure framework**.

## What Optivor Is

- **Runtime Engine**: Performs high-performance image transformations (resize, format conversion to WebP/AVIF, quality optimization) using `libvips`.
- **Bring-Your-Own-Storage (BYOS)**: Integrates directly with object storage you already own (AWS S3, MinIO, Cloudflare R2, Backblaze B2, Google Cloud Storage).
- **Extensible Framework**: Offers out-of-process driver conventions and deployment adapters.
- **Production-Grade Infrastructure**: Includes built-in rate limiting, LRU filesystem caching, OpenTelemetry tracing, and Prometheus metrics.

## What Optivor Is NOT

- **Not an Image Host**: Optivor does not store your original images.
- **Not a CDN**: Optivor runs behind your CDN or directly as an origin microservice.
- **Not SaaS/Open-Core**: There are no paid tiers or enterprise-only binary blobs.

## Target Audience

Engineering teams building scalable web applications who want full control over their image pipelines without vendor lock-in or recurring SaaS per-image transformation costs.
