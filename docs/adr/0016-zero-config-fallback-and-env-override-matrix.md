# ADR-0016: Zero-Config Fallback & Environment Variable Override Matrix

- **Status**: Accepted
- **Deciders**: Optivor Core Architecture Team
- **Date**: 2026-07-25

---

## Context and Problem Statement

To minimize developer onboarding friction and support seamless 1-click cloud deployments (e.g. Railway, Render, Fly.io), Optivor must support zero-config local execution and 100% environment variable configuration. 

Previously, starting Optivor without a valid `optivor.yaml` containing S3 endpoint credentials resulted in validation errors. Developers testing Optivor locally had to scaffold configuration and run a local S3/MinIO container even for simple trial setups. Additionally, cloud platform deployments often supply configuration solely via environment variables rather than mounted YAML files.

---

## Decision Outcome

We decide to implement **Zero-Config Local Storage Fallback** and a formalized **Environment Variable Precedence Hierarchy**:

1. **Zero-Config Local Fallback**:
   - If no `optivor.yaml` is provided and no S3 environment variables are set, `optivor start` automatically initializes a local storage driver pointing to `./storage` (or `/tmp/optivor-storage`).
   - The server serves and optimizes images directly from local disk out-of-the-box without requiring S3 credentials or configuration files.

2. **Environment Variable Override Hierarchy**:
   - Configuration priority order: `Environment Variables` > `YAML Configuration File` > `Default Values`.
   - All nested YAML keys can be overridden via `OPTIVOR_<SECTION>_<KEY>` (e.g., `OPTIVOR_SERVER_PORT=8080`, `OPTIVOR_CACHE_TYPE=redis`).
   - Shorthand cloud environment variables (`OPTIVOR_S3_BUCKET`, `OPTIVOR_S3_ENDPOINT`, `OPTIVOR_S3_REGION`, `OPTIVOR_S3_ACCESS_KEY_ID`, `OPTIVOR_S3_SECRET_ACCESS_KEY`) are automatically normalized into `storage.s3.*`.

3. **Interactive CLI Scaffolding**:
   - `optivor init --interactive` prompts developers step-by-step for provider credentials, tests connection live, and generates standard `optivor.yaml`.

---

## Precedence & Configuration Matrix

| Priority | Source | Example | Description |
| :--- | :--- | :--- | :--- |
| 1 (Highest) | System Env Vars | `OPTIVOR_STORAGE_S3_BUCKET=prod-images` | Direct environment override |
| 2 | Shorthand Env Vars | `OPTIVOR_S3_BUCKET=prod-images` | PaaS/Serverless cloud shorthand |
| 3 | YAML Config File | `storage.s3.bucket: dev-images` | `optivor.yaml` file configuration |
| 4 (Lowest) | Defaults | `local` driver / `./storage` | Built-in zero-config fallback |

---

## Consequences

### Positive
- Developers can test Optivor locally in under 5 seconds with zero setup.
- Containerized deployments on Railway, Render, Fly.io work without mounting config files.
- Clear, predictable configuration precedence across all environments.

### Negative
- Requires maintaining environment variable mapping aliases in the configuration loader.
