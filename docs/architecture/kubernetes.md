# Architecture: Kubernetes & Helm Deployment Adapter

## Overview

Per **ADR-0014**, the Kubernetes and Helm deployment adapter layer is designed to package, secure, and scale Optivor instances in multi-tenant, high-availability container orchestration environments.

## Architectural Mechanisms

### 1. Zero-Trust Security Configuration
To satisfy strict corporate security standards, the Helm chart enforces a zero-trust execution model:
*   **ReadOnly Root Filesystem:** Prevents malware or runtime drift from modifying container binary assets.
*   **Sandboxed Temporary Memory:** Since libvips requires a writable directory to cache image transposition buffers, an `emptyDir` memory-backed volume is mounted to `/tmp`.
*   **Minimal Privileges:** Runs as UID/GID `10001` with `runAsNonRoot: true` and all Linux capabilities explicitly dropped (`drop: [ALL]`).

### 2. ConfigMap Checksum Hashing & Rolling Restarts
To ensure configuration updates inside `values.yaml` take effect immediately without requiring manual CLI restarts:
*   A SHA-256 hash of the rendered `optivor.yaml` ConfigMap is computed and injected as an annotation on the Pod template:
    `checksum/config: <sha256sum>`
*   When a user updates their config via `helm upgrade`, the annotation changes, forcing Kubernetes to perform a safe rolling restart of all Optivor replicas.

### 3. Stateless Cache Scaling (Redis Integration)
*   **Multi-Replica Sync:** With multiple running replicas (controlled by replicaCount or HPA), the default filesystem cache becomes inconsistent.
*   The deployment architecture integrates a centralized **Redis Cache Adapter** where all replicas write and read image payloads concurrently, eliminating duplicate libvips processing spikes across different pods.
