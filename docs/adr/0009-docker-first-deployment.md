# ADR-0009: Docker-First Deployment Strategy

## Status

Accepted

## Context

Optivor initially targeted standalone binary execution and systemd service deployment (ADR-0006). As Optivor moves towards production deployments and provider-agnostic infrastructure (V0.5), containerization via Docker is required as a primary distribution mechanism.

Deployment environments vary widely, ranging from single-node Docker hosts and docker-compose setups to managed container platforms. However, core architecture principles (ADR-0001, ADR-0002) dictate zero platform lock-in and strict separation between the runtime core and environment configuration.

## Decision

1. **Docker as Primary Distribution Mechanism**:
   Docker container image (`optivor:latest`) is established as the primary delivery format alongside standalone binaries. The existing systemd deployment mechanism (`make install`) remains supported and un-deprecated.

2. **Configuration Injection via Volume Mount & Environment Overrides**:
   - The primary configuration file `optivor.yaml` is provided to the container via volume mount:
     `-v ./optivor.yaml:/etc/optivor/optivor.yaml:ro`
   - Sensitive credentials (e.g. S3 secrets, HMAC secrets) can be injected via standard environment variable overrides (e.g., `OPTIVOR_STORAGE_S3_SECRET_ACCESS_KEY`, `OPTIVOR_AUTH_SECRET`).

3. **Runtime `--provider` Startup Flag**:
   - Introduce a `--provider <name>` startup flag (e.g. `optivor --provider r2`, `optivor --provider minio`).
   - The `--provider` flag dynamically overrides the `storage.driver` field in `optivor.yaml` without requiring file edits.
   - If an unknown provider name is supplied at startup, the process terminates immediately with `exit 1` and a clear error message.

4. **Container Health Checking**:
   - Implement an explicit `HEALTHCHECK` directive in the official `Dockerfile` targeting the `/healthz` endpoint (30s interval, 5s timeout, 3 retries).

5. **Reference Compose Template**:
   - Include a standard, production-ready `docker-compose.yml` in the project root showcasing volume mounting, environment secret overrides, container healthchecks, and local MinIO development service.

## Consequences

- Optivor can be deployed seamlessly in any containerized environment (Docker Compose, ECS, Kubernetes, Nomad).
- Zero file mutation required when deploying across development, staging, and production thanks to `--provider` and environment override mechanisms.
- Systemd standalone users retain full compatibility without breaking changes.
