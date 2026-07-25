# PaaS Deployment Blueprint: Railway, Render & Fly.io

This guide provides 1-click container deployment configurations and instructions for deploying Optivor on modern PaaS platforms using 100% environment variable configuration.

---

## 1. Railway Deployment

Optivor can be deployed directly to Railway using standard environment variables without mounting a configuration file.

### Environment Variables
Set the following environment variables in your Railway project service settings:

```env
PORT=8080
OPTIVOR_S3_ENDPOINT=https://<account-id>.r2.cloudflarestorage.com
OPTIVOR_S3_BUCKET=my-app-assets
OPTIVOR_S3_REGION=auto
OPTIVOR_S3_ACCESS_KEY_ID=your-access-key
OPTIVOR_S3_SECRET_ACCESS_KEY=your-secret-key
OPTIVOR_CACHE_TYPE=redis
OPTIVOR_CACHE_REDIS_ADDR=redis.railway.internal:6379
```

---

## 2. Render Deployment

Deploy Optivor on Render using a Web Service running the official `ghcr.io/optivor/optivor:latest` container image.

### `render.yaml` Blueprint

```yaml
services:
  - type: web
    name: optivor-engine
    env: docker
    imageUrl: ghcr.io/optivor/optivor:latest
    plan: starter
    region: oregon
    envVars:
      - key: PORT
        value: 8080
      - key: OPTIVOR_S3_ENDPOINT
        value: https://s3.us-east-1.amazonaws.com
      - key: OPTIVOR_S3_BUCKET
        value: prod-image-bucket
      - key: OPTIVOR_S3_REGION
        value: us-east-1
      - key: OPTIVOR_S3_ACCESS_KEY_ID
        sync: false
      - key: OPTIVOR_S3_SECRET_ACCESS_KEY
        sync: false
```

---

## 3. Fly.io Deployment

Deploy Optivor as a ultra-fast edge container on Fly.io across global regions.

### `fly.toml` Configuration

```toml
app = "optivor-engine"
primary_region = "iad"

[build]
  image = "ghcr.io/optivor/optivor:latest"

[env]
  PORT = "8080"
  OPTIVOR_CACHE_TYPE = "fs"
  OPTIVOR_CACHE_FS_DIR = "/tmp/optivor-cache"

[[services]]
  internal_port = 8080
  protocol = "tcp"

  [[services.ports]]
    handlers = ["http"]
    port = 80

  [[services.ports]]
    handlers = ["tls", "http"]
    port = 443
```

Set secrets via Fly CLI:

```bash
fly secrets set \
  OPTIVOR_S3_ENDPOINT="https://s3.us-east-1.amazonaws.com" \
  OPTIVOR_S3_BUCKET="my-fly-bucket" \
  OPTIVOR_S3_ACCESS_KEY_ID="my-key" \
  OPTIVOR_S3_SECRET_ACCESS_KEY="my-secret"
```
