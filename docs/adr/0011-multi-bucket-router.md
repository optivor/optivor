# ADR-0011: Multi-Bucket Router and Per-Bucket Access Policy

- **Status:** Accepted
- **Date:** 2026-07-20
- **Authors:** Optivor Team

## Context

Prior to V0.6, Optivor was configured with a single default storage bucket (e.g. AWS S3). As workloads scale across multi-cloud infrastructure, applications require managing multiple storage buckets across distinct providers (S3, Cloudflare R2, Backblaze B2, Google Cloud Storage) from a single runtime instance.

Additionally, different asset types demand distinct authorization mechanisms:
- Public media assets (`access: public`) require no authorization signatures.
- Sensitive or restricted media (`access: signed`) require HMAC-SHA256 URL signatures.
- Internal-only assets (`access: private`) must be blocked from external public access.

Furthermore, storage providers experience regional outages, requiring automatic cross-provider failover chains.

## Decision

1. **Multi-Bucket Declarative Schema (`optivor.yaml`)**:
   Introduce a `buckets` list in `optivor.yaml`. Each entry defines:
   - `name`: Bucket URL alias (e.g., `primary-images`, `secure-assets`).
   - `provider`: Storage driver type (`s3`, `r2`, `b2`, etc.).
   - `bucket`, `region`, `account_id`, `endpoint`: Provider configuration parameters.
   - `access_key_id`, `secret_access_key`: Credentials (overridable via `OPTIVOR_BUCKET_<NAME>_*`).
   - `access`: Access policy (`public`, `signed`, `private`).
   - `fallback`: Alias of backup bucket for failover (optional).

2. **URL Routing Scheme**:
   Media URLs follow the pattern `/image/<bucket-alias>/<object-key>?w=...`.
   Legacy single-bucket URLs (`/image/<object-key>`) are routed to the default bucket if configured for backward compatibility.

3. **Bucket Router Interface (`internal/storage/router`)**:
   - `Resolve(ctx, alias)`: Locates the primary driver or falls back to the configured backup bucket upon storage error.
   - `Policy(alias)`: Returns the `AccessPolicy` enum (`PolicyPublic`, `PolicySigned`, `PolicyPrivate`).
   - Circular failover chains (A -> B -> A) are detected and rejected at startup validation.

4. **Layer Boundaries (ADR-0002 Compliance)**:
   The `router` package sits in `internal/storage/router`. The HTTP server layer calls `BucketRouter` to check authorization policy and obtain the storage driver without coupling to provider specifics.

## Consequences

- **Positive**: Single Optivor deployment handles multiple buckets across different providers with per-bucket security policies and resilient failover.
- **Backward Compatibility**: Single-bucket configuration in `optivor.yaml` continues to work cleanly.
