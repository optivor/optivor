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

---

> [!NOTE]
> **Multi-Bucket Routing & Access Control:** Optivor V0.6+ supports multi-bucket declarative routing (`buckets[]`) with per-bucket security policies (`public`, `signed`, `private`) and cross-provider failover chains.
>
> [!NOTE]
> **Signed URLs & Authentication:** HMAC-SHA256 URL signing (`auth.signed_urls.enabled: true`) enforces `sig` and `expires` signatures.
>
> [!IMPORTANT]
> **DoS Protection Notice:** Enforce reasonable `max_width` and `max_height` values in `optivor.yaml` to prevent decompression-bomb attacks.

---

## 5-Minute Quick Start

### 1. Build the Binary

Ensure `libvips` development header is installed on your Linux system (`apt install libvips-dev`):

```bash
make build
```

This compiles the standalone Go runtime binary to `bin/optivor`.

### 2. Configure `optivor.yaml`

Initialize a new `optivor.yaml` configuration file:

```bash
./bin/optivor init
```

Or configure multi-bucket storage routing:

```yaml
server:
  port: 8080
  read_timeout: 30s
  write_timeout: 30s
  request_timeout: 30s
  log_level: "info"      # debug | info | warn | error
  log_format: "text"     # text | json
  rate_limit:
    enabled: true
    rps: 10
    burst: 20
  image:
    max_width: 5000   # px DoS limit
    max_height: 5000  # px DoS limit

buckets:
  - name: "primary-images"
    provider: s3
    endpoint: "https://s3.us-east-1.amazonaws.com"
    bucket: "my-aws-bucket"
    region: "us-east-1"
    access: public

  - name: "secure-assets"
    provider: r2
    endpoint: "https://account-id.r2.cloudflarestorage.com"
    bucket: "my-r2-bucket"
    access: signed
    fallback: "primary-images"

cache:
  fs:
    dir: "/tmp/optivor-cache"
    max_size_mb: 1024

telemetry:
  enabled: true
  otlp_endpoint: ""
  service_name: "optivor"
  sampling_ratio: 1.0

image:
  contain_background_color: "#ffffff"
```

### 3. Run Optivor

```bash
./bin/optivor -config ./optivor.yaml
```

### 4. Serve Resized & WebP Converted Images

Request images via single-bucket key or multi-bucket alias:

```bash
# Multi-bucket alias route
curl -i "http://localhost:8080/image/primary-images/products/123/main.jpg?w=300&h=300&fit=cover&format=webp"
```

Response returns `Content-Type: image/webp` and `X-Optivor-Cache: MISS` (or `HIT` on subsequent requests).

### 5. Health Check & Diagnostics

Check system health using `/healthz`, `/health`, or `/healtz` endpoints or running `optivor doctor`:

```bash
curl -i "http://localhost:8080/healthz"
./bin/optivor doctor
```

### 6. Prometheus Metrics

Optivor exposes Prometheus metrics at `GET /metrics`:

```bash
curl -i "http://localhost:8080/metrics"
```

### 7. Storage Driver Management via CLI

Install storage provider drivers from local paths, GitHub repository shorthands, or direct release URLs:

```bash
# Install driver via GitHub repository shorthand
./bin/optivor driver install github:optivor/optivor-driver-r2@v1.2.0

# Install driver via direct HTTPS release URL
./bin/optivor driver install https://github.com/optivor/optivor-driver-r2/releases/download/v1.2.0/optivor-driver-r2-linux-amd64

# List registered drivers
./bin/optivor driver list
```

### 8. Bucket Lifecycle Management via CLI

Manage multi-cloud bucket retention rules:

```bash
./bin/optivor bucket lifecycle list primary-images
./bin/optivor bucket lifecycle set primary-images --ttl-days 30
./bin/optivor bucket lifecycle delete primary-images --all
```

---

## Core Principles

- **Bring Your Own Storage.** Optivor never stores your images. Your data lives in your bucket, under your account, under your control — always.
- **Open-source-first, not open-core.** There is no paid tier hiding behind the free one. What you see in this repository is the product.
- **Provider-agnostic by default.** The runtime works on a plain VM with nothing but a binary and a config file.
- **Composable, not monolithic.** Storage, transformation, caching, and deployment are independently replaceable.
- **Contributor-first.** You can own an entire piece of this project — a storage driver, a deployment adapter — without needing to understand the whole codebase.

---

## Architecture Overview

```
CLI
  ↓
Runtime (Router & Pipeline)
  ↓
Storage Drivers (S3, R2, B2, GCS)
  ↓
Object Storage (yours)
```

The runtime knows nothing about cloud providers. Deployment adapters know nothing about image processing. Storage drivers know nothing about either. Each piece does one job.

Full reasoning lives in [`docs/adr/`](./docs/adr).

---

## Documentation

Explore the complete Optivor documentation in [`docs/wiki/`](./docs/wiki):

- [Introduction](./docs/wiki/introduction.md) — Framework overview and philosophy
- [Quick Start Guide](./docs/wiki/quick-start.md) — Setup with Docker and binary
- [CLI Reference](./docs/wiki/cli-reference.md) — Command and flag reference
- [Configuration Reference](./docs/wiki/configuration.md) — `optivor.yaml` schema & environment overrides
- [Multi-Cloud & Multi-Bucket Management](./docs/wiki/multi-cloud-management.md) — Multi-bucket configuration & security policies
- [Edge Integration Guide](./docs/wiki/edge-integration.md) — CDN and Cloudflare Workers setup
- [Storage Driver Guide](./docs/wiki/storage-drivers.md) — Building & installing custom storage drivers
- [Storage Driver SDK Specification](./docs/wiki/driver-sdk-specification.md) — Out-of-process IPC protocol specification
- [FAQ](./docs/wiki/faq.md) — Frequently asked questions

---

## Roadmap

See [`ROADMAP.md`](./ROADMAP.md) for completed milestones and future trajectory.

## License

Apache License 2.0. See [`LICENSE`](./LICENSE).
