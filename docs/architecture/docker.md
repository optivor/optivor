# Architecture: Docker-First Deployment

## Overview

Per **ADR-0009** (Docker-First Deployment Strategy), Docker containerization is established as a primary delivery model alongside standalone binaries.

## Key Design Mechanisms

1. **Volume-Mounted Configuration**:
   - `optivor.yaml` is mounted read-only to `/etc/optivor/optivor.yaml`.
   - Immutable binary inside container reads config directly from filesystem mount.

2. **Environment Variable Injection**:
   - Secrets are overridden at container startup using `OPTIVOR_*` environment variables.

3. **Dynamic `--provider` Flag**:
   - Container startup command can specify `--provider <name>` to override storage backend without editing configuration files.

4. **Container Health Checking**:
   - Built-in `HEALTHCHECK` directive queries `/healthz` to inform orchestrators (Docker Swarm, Kubernetes, ECS) of container health status.
