# ADR-0010: Provider Registry & Driver Binary Convention

## Status

Accepted

## Context

As Optivor expands support across cloud platforms (Cloudflare R2, Backblaze B2, Google Cloud Storage, AWS S3, local MinIO), maintaining vendor-specific SDKs and API drivers inside the core repository creates maintenance overhead, binary bloat, and violates the provider-agnostic core architecture (ADR-0001, ADR-0002).

Per **ADR-0003** (Extensibility Model), external integrations and adapters must run out-of-process to protect core runtime stability and isolation.

## Decision

1. **Provider-Agnostic Core Repository**:
   - The core Optivor repository contains **no provider-specific code** beyond the universal S3-compatible driver (`internal/storage/s3`), which acts as the default baseline protocol.
   - Provider-specific experiences (e.g. Cloudflare API integration, B2 lifecycle rules, GCS IAM authentication) ship in separate repositories following the `optivor-driver-<name>` naming convention (e.g. `optivor-driver-r2`, `optivor-driver-b2`).

2. **Standalone Binary Drivers (Terraform Pattern)**:
   - External storage drivers build as standalone executable binaries.
   - Core Optivor communicates with registered driver binaries out-of-process via standard subprocess execution.

3. **Handshake & Metadata Specification**:
   - Every driver binary must respond to the `--optivor-handshake` flag by printing JSON metadata to stdout and exiting with status 0:
     ```json
     {
       "name": "r2",
       "version": "0.1.0",
       "optivor_api": "v1"
     }
     ```

4. **Local Driver Registry**:
   - The CLI manages installed drivers in a local registry file stored at `~/.config/optivor/drivers.json`.
   - The `optivor driver` CLI subcommands (`install`, `list`, `remove`, `info`) validate handshakes and manage entries in `drivers.json`.

## Consequences

- Core Optivor remains lightweight, decoupled, and free of vendor SDK bloat.
- Community and commercial third parties can independently build, release, and distribute custom storage drivers.
- Driver execution failure will not crash the primary Optivor process.
