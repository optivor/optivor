# Zero-Config Quickstart Guide

Get Optivor running locally on your computer in under 10 seconds without cloud credentials, S3 buckets, or YAML files.

---

## 1. Instant Zero-Config Storage Mode

When Optivor starts without an `optivor.yaml` file or environment variables, it automatically boots into **Zero-Config Local Storage Fallback Mode**:

```bash
# Start Optivor immediately
./bin/optivor
```

Output:
```text
time=2026-07-25T23:00:00 level=INFO msg="Using zero-config local storage fallback" directory="./storage"
time=2026-07-25T23:00:00 level=INFO msg="Optivor binary started successfully" port=8080
```

---

## 2. Serving Local Images

Place any test image inside the `./storage` directory:

```bash
mkdir -p ./storage
cp sample.jpg ./storage/sample.jpg
```

Request optimized images via HTTP:

```bash
# Convert to WebP with width 400
curl "http://localhost:8080/sample.jpg?w=400&f=webp" --output output.webp

# Convert to AVIF with smart entropy crop
curl "http://localhost:8080/sample.jpg?w=300&h=300&fit=smart&f=avif" --output avatar.avif
```

---

## 3. Environment Variable Quick-Start

If you want to quickly connect to an existing S3 or Cloudflare R2 bucket without creating a `yaml` file:

```bash
export OPTIVOR_S3_ENDPOINT="https://s3.us-east-1.amazonaws.com"
export OPTIVOR_S3_BUCKET="my-test-bucket"
export OPTIVOR_S3_ACCESS_KEY_ID="AKIA..."
export OPTIVOR_S3_SECRET_ACCESS_KEY="secret..."

./bin/optivor
```
