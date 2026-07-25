# PaaS Deployment Blueprint: Railway, Render & Fly.io

This guide provides step-by-step instructions for deploying Optivor on modern PaaS platforms (Railway, Render, Fly.io) using either **Direct GitHub Repository Builds (Zero GHCR setup)** or **GitHub Container Registry (GHCR) images**.

---

## Deployment Strategies

### Option A: Direct GitHub Repository Build (Recommended - Zero Registry Required)
Both Railway, Render, and Fly.io can build directly from your GitHub repository using the `Dockerfile` included in the root directory. **You do not need to build or push container images to GHCR beforehand.**

### Option B: GHCR Container Image Deployment
If you prefer deploying pre-built container images, GitHub Actions automatically builds and publishes `ghcr.io/optivor/optivor:latest` whenever a release tag is pushed.

---

## 1. Railway Deployment

### Method A: Deploy Direct from GitHub Repo (No GHCR required)
1. Log into [Railway.app](https://railway.app).
2. Click **New Project** -> **Deploy from GitHub repo**.
3. Select your `optivor` repository.
4. Railway automatically detects the root `Dockerfile` and builds the project.
5. In **Variables**, add your environment variables:

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

### Method A: Deploy Direct from GitHub Repo (No GHCR required)
1. Log into [Render.com](https://render.com).
2. Click **New +** -> **Web Service**.
3. Connect your GitHub repository (`optivor`).
4. Set Environment to **Docker** (leave Build Command and Start Command blank as Render reads `Dockerfile`).
5. Add required environment variables under **Environment Variables**.

### Method B: `render.yaml` Infrastructure-as-Code

```yaml
services:
  - type: web
    name: optivor-engine
    env: docker
    # Option B: Image URL (or comment out imageUrl if deploying from repo)
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

### Method A: Deploy Direct from Local Source / GitHub (No GHCR required)
Run Fly CLI in your repository root:

```bash
fly launch
```

Fly automatically detects `Dockerfile`, builds it remotely, and deploys the container.

### `fly.toml` Configuration

```toml
app = "optivor-engine"
primary_region = "iad"

[build]
  # Leave dockerfile specified if building from repo, or set image for GHCR
  dockerfile = "Dockerfile"

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
