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

> [!WARNING]
> **V0 Security Notice:** Optivor V0 does **not** include built-in request signing or authentication (deferred to V0.1). Do **not** expose Optivor directly to the public internet without an external authentication proxy or CDN protection layer.
>
> [!CAUTION]
> **V0 Cache Growth Notice:** The filesystem cache in V0 grows continuously without automated LRU eviction (deferred to V0.1). Manage disk space for your cache directory (`/tmp/optivor-cache`) manually or via cron tasks in production.
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

Create an `optivor.yaml` file in your current working directory:

```yaml
server:
  port: 8080
  read_timeout: 30s
  write_timeout: 30s
  image:
    max_width: 5000   # px DoS limit
    max_height: 5000  # px DoS limit

storage:
  s3:
    endpoint: "https://s3.amazonaws.com"
    bucket: "my-image-bucket"
    region: "us-east-1"
    access_key_id: "YOUR_ACCESS_KEY"
    secret_access_key: "YOUR_SECRET_KEY"

cache:
  fs:
    dir: "/tmp/optivor-cache"

image:
  contain_background_color: "#ffffff"
```

### 3. Run Optivor

```bash
./bin/optivor -config ./optivor.yaml
```

### 4. Serve Resized & WebP Converted Images

Request any S3 object key with image transformation parameters:

```bash
curl -i "http://localhost:8080/image/products/123/main.jpg?w=300&h=300&fit=cover&format=webp"
```

Response will return with `Content-Type: image/webp` and header `X-Optivor-Cache: MISS` (or `HIT` on subsequent requests).

---

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

---

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

The runtime knows nothing about cloud providers. Deployment adapters
know nothing about image processing. Storage drivers know nothing about
either. Each piece does one job.

The full reasoning behind these boundaries lives in [`docs/adr/`](./docs/adr).

---

## Roadmap

See [`ROADMAP.md`](./ROADMAP.md) for current milestone status and features deferred to post-V0 releases.

## License

Apache License 2.0. See [`LICENSE`](./LICENSE).
