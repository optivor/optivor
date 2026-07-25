# ADR-0014: Redis Cache Backend, Smart Cropping, and Helm Deployment Adapter

- **Status:** Accepted
- **Date:** 2026-07-25
- **Authors:** Optivor Team

## Context

As Optivor transitions to a production-grade infrastructure component (V1), it requires support for enterprise deployment primitives, multi-node horizontal scaling, and modern image optimization features. Specifically:
1. File-system based caching binds container instances to local disks, preventing horizontal scaling in elastic environments (like Kubernetes pods) without complex volume mount synchronization.
2. Manual resize modes can clip critical subjects in images, requiring smart/content-aware focal cropping.
3. Container deployments need strict security policies (read-only filesystems, non-root execution) and structured configuration maps to comply with enterprise orchestration guidelines.

## Decision

1. **Redis Cache Backend (`internal/cache/redis`)**:
   - Implement the `cache.Cache` interface using `github.com/redis/go-redis/v9`.
   - Pack the binary payload and Content-Type metadata into a single serialized record to retrieve cached images in one network round-trip.
   - Allow horizontal scaling of stateless Optivor runtime instances backed by a centralized Redis cluster.

2. **Smart Cropping (`fit=smart`)**:
   - Integrate `vips.InterestingAttention` (entropy-based visual attention mapping) in the core `internal/pipeline` package.
   - Expose the `fit=smart` query parameter to let clients auto-crop around the most visually significant parts of the image (e.g., human faces or high-contrast foreground objects).

3. **Helm Chart Deployment Adapter (`deploy/helm/optivor`)**:
   - Create a customizable Helm chart providing out-of-the-box Kubernetes manifests.
   - Enforce secure defaults: `runAsNonRoot: true` (user `10001`), `readOnlyRootFilesystem: true`, and dropped capabilities (`drop: [ALL]`).
   - Support horizontal autoscaling (HPA) based on CPU utilization and mount configurations via standard ConfigMaps and Secrets.

## Consequences

- **Scalability:** Optivor pods can be scaled dynamically without losing cache state or triggering cache synchronization overhead.
- **Visual Quality:** Users gain access to content-aware cropping, reducing crop-clipping mistakes on custom layouts.
- **Ops Readiness:** Simplified deployment and management inside Kubernetes clusters with standard Helm charts conforming to cloud-native security practices.
