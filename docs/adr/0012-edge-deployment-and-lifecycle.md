# ADR-0012: Edge Deployment Strategy and Lifecycle Management API

- **Status:** Accepted
- **Date:** 2026-07-20
- **Authors:** Optivor Team

## Context

As media architectures evolve, applications require edge integration models (Cloudflare Workers, AWS Lambda@Edge) and multi-cloud bucket lifecycle rule automation (S3, R2, B2, GCS expiration/transition policies).

## Decision

1. **Edge Deployment Models**:
   - **Proxy Mode (Primary)**: Cloudflare Workers or Lambda@Edge route image requests to the Optivor origin instance while providing CDN edge caching.
   - **Edge Transform Mode (Alternative)**: Lightweight edge workers process basic transforms directly at the edge when full vips acceleration is unnecessary.

2. **Bucket Lifecycle CLI API**:
   - `optivor bucket lifecycle list <alias>`
   - `optivor bucket lifecycle set <alias> --ttl-days 30`
   - `optivor bucket lifecycle set <alias> --rule-file lifecycle.yaml`
   - `optivor bucket lifecycle delete <alias> [--rule-id id | --all]`

3. **Provider Translation**:
   Lifecycle configurations are parsed from a provider-agnostic YAML definition (`lifecycle.yaml`) and mapped to the target storage driver API.

## Consequences

- Standardized CLI operational workflow for lifecycle automation across multi-cloud buckets.
- Clear pattern for edge CDN integrations.
