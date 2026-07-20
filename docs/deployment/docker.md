# Docker Deployment Guide

Optivor provides first-class support for containerized deployment via Docker and Docker Compose.

## 1. Quick Start with Docker

Build the official container image:

```bash
make docker-build
# or
docker build -t optivor:latest .
```

Run Optivor with a volume-mounted configuration:

```bash
docker run -d \
  --name optivor \
  -p 8080:8080 \
  -v $(pwd)/optivor.yaml:/etc/optivor/optivor.yaml:ro \
  optivor:latest
```

Verify health status:

```bash
docker inspect --format='{{json .State.Health}}' optivor
```

## 2. Environment Variable Overrides

Sensitive variables or runtime parameters can be injected directly via environment variables:

```bash
docker run -d \
  --name optivor \
  -p 8080:8080 \
  -v $(pwd)/optivor.yaml:/etc/optivor/optivor.yaml:ro \
  -e OPTIVOR_STORAGE_S3_SECRET_ACCESS_KEY="my-secret-key" \
  -e OPTIVOR_AUTH_SECRET="my-hmac-secret" \
  optivor:latest
```

## 3. Dynamic Storage Provider Selection (`--provider`)

Override the storage driver at container startup using the `--provider` flag:

```bash
docker run -d \
  --name optivor \
  -p 8080:8080 \
  -v $(pwd)/optivor.yaml:/etc/optivor/optivor.yaml:ro \
  optivor:latest --provider r2
```

Supported built-in providers include `s3`, `minio`, and `r2`. Providing an unsupported provider name will exit immediately with an error.

## 4. Docker Compose Setup

Run Optivor alongside a local MinIO S3-compatible service:

```bash
docker-compose up -d
```

The reference `docker-compose.yml` mounts `./optivor.yaml.example` to `/etc/optivor/optivor.yaml` inside the container and configures standard health checks.
