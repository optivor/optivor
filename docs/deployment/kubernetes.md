# Kubernetes & Helm Deployment Guide

Optivor is designed for cloud-native environments and includes a production-ready Helm chart supporting high-availability, horizontal scalability, and non-root execution policies.

## 1. Prerequisites

- A Kubernetes cluster (v1.20+)
- Helm 3 installed
- A S3-compatible storage bucket (AWS S3, Cloudflare R2, MinIO, etc.)
- A Redis instance (optional, required if using horizontal scaling cache)

## 2. Deploying with Helm

The Helm chart is located at `deploy/helm/optivor`.

### Quick Installation

Install the chart using the default filesystem cache (`fs`):

```bash
helm install optivor ./deploy/helm/optivor \
  --set config.storage.s3.endpoint="https://s3.amazonaws.com" \
  --set config.storage.s3.bucket="my-optivor-bucket"
```

### High Availability with Redis Cache

For horizontal scaling (multiple replicas), configure the Redis cache backend. This allows multiple Optivor pods to share cache state seamlessly:

```bash
helm install optivor ./deploy/helm/optivor \
  --set replicaCount=3 \
  --set config.cache.type="redis" \
  --set config.cache.redis.addr="my-redis-host:6379" \
  --set config.storage.s3.endpoint="https://s3.amazonaws.com" \
  --set config.storage.s3.bucket="my-optivor-bucket"
```

---

## 3. Configuration & Values Reference

Key configurable values in `deploy/helm/optivor/values.yaml` include:

| Parameter | Description | Default |
| :--- | :--- | :--- |
| `replicaCount` | Number of running pods | `2` |
| `image.repository` | Container image repository | `optivor/optivor` |
| `resources.limits.cpu` | CPU limit for libvips execution threads | `2` |
| `resources.limits.memory` | Memory limit for processing pipelines | `2Gi` |
| `config.cache.type` | Caching backend type (`fs` or `redis`) | `fs` |
| `config.cache.redis.addr` | Host and port of the Redis cache service | `redis-master:6379` |
| `config.storage.driver` | Active storage provider backend | `s3` |
| `config.storage.s3.endpoint` | Target S3 endpoint URL | `http://optivor-minio:9000` |
| `config.storage.s3.bucket` | Default image storage bucket | `optivor-images` |

---

## 4. Production Security Controls

The Helm chart enforces strict security policies by default to comply with enterprise cloud environments:

*   **Non-Root Execution:** Containers run as UID `10001` (`runAsNonRoot: true`).
*   **ReadOnly Root Filesystem:** Root FS is mounted read-only. Pods write temporary files to a secure `/tmp` memory-backed `emptyDir` sandbox.
*   **Privilege Controls:** Privilege escalation is explicitly disabled (`allowPrivilegeEscalation: false`).
*   **Capability Drops:** All Linux capabilities are dropped (`drop: [ALL]`).

---

## 5. Secret Reference & Credentials Injection

To avoid placing sensitive credentials (like S3 credentials or Redis passwords) in `values.yaml`, specify an existing Kubernetes secret containing the following keys using the `existingSecret` setting:

```yaml
# values.yaml snippet
existingSecret: "optivor-production-secrets"
```

The secret `optivor-production-secrets` should contain the following fields:

*   `OPTIVOR_STORAGE_S3_ACCESS_KEY_ID`: Access key ID for the global S3 storage bucket.
*   `OPTIVOR_STORAGE_S3_SECRET_ACCESS_KEY`: Secret access key for the global S3 storage bucket.
*   `OPTIVOR_CACHE_REDIS_PASSWORD`: Optional Redis password.
*   `OPTIVOR_AUTH_SIGNED_URLS_SECRET`: HMAC secret (if URL signature verification is enabled).
