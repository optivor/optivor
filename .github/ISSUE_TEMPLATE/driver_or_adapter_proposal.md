---
name: "🔌 Storage Driver or Adapter Proposal"
about: Propose a new Storage Driver (e.g. S3, GCS, Azure, Cloudflare R2, Ceph) or Cache Adapter
title: "[Proposal]: "
labels: proposal, driver
assignees: ''
---

## 🎯 Proposal Summary
A clear summary of the storage provider or caching adapter being proposed (e.g., *Ceph RADOS Object Storage Driver*, *Cloudflare R2 Native Binding*, *NATS JetStream Cache Backend*).

## 🏢 Target Storage Provider / Backend Specs
- **Target Technology**: (e.g. MinIO, Cloudflare R2, Backblaze B2, Azure Blob Storage)
- **Protocol / SDK**: (e.g. S3-compatible REST API, Native Go SDK)
- **Authentication Method**: (Access/Secret Key, OAuth2, IAM Role, IAM Service Account)

## 📋 Conformance & Interface Checklist
Optivor storage drivers must satisfy the [`storage.Driver`](./internal/storage/driver.go) interface contract:
- [ ] `Get(ctx, key)` -> Returns `io.ReadCloser` or `ErrNotFound`
- [ ] `Put(ctx, key, reader, size, contentType)` -> Streaming upload
- [ ] `Delete(ctx, key)` -> Object deletion
- [ ] `Exists(ctx, key)` -> Fast metadata HEAD check

## 🧪 Testing Plan & Mocking Strategy
How will unit and integration tests be executed for this driver? (e.g., Docker container minio, testcontainers-go, local emulator).

## 📄 References & Specs
Link to official documentation or SDK specifications for the underlying storage engine.
