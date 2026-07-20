# Quick Start Guide

Get up and running with Optivor in under 5 minutes using Docker.

## Step 1: Clone and Prepare Config

Clone the repository and copy the example configuration:

```bash
git clone https://github.com/optivor/optivor.git
cd optivor
cp optivor.yaml.example optivor.yaml
```

Update `optivor.yaml` with your storage credentials (e.g. S3 endpoint, bucket, access keys).

## Step 2: Run with Docker Compose

Start Optivor alongside a local MinIO service for testing:

```bash
docker-compose up -d
```

## Step 3: Run standalone with Docker

Or run the container with a volume-mounted configuration:

```bash
docker run -d \
  --name optivor \
  -p 8080:8080 \
  -v $(pwd)/optivor.yaml:/etc/optivor/optivor.yaml:ro \
  optivor:latest
```

## Step 4: Request an Image Transformation

Transform an image on the fly:

```bash
curl -i "http://localhost:8080/image/sample.jpg?w=400&h=300&fit=cover&format=webp"
```

Response will include transformed WebP binary payload and `X-Optivor-Cache: MISS` header.
