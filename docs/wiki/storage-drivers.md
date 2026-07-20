# Storage Driver Development Guide

Optivor enforces a provider-agnostic core architecture (ADR-0002, ADR-0010). Storage drivers for cloud providers live outside the core repository as standalone binaries following the `optivor-driver-<name>` naming convention.

## Driver Specification

An Optivor storage driver is an executable binary that satisfies the out-of-process protocol interface.

### 1. Handshake Protocol (`--optivor-handshake`)

When invoked with `--optivor-handshake`, the binary must print a JSON object to standard output and exit with code 0:

```json
{
  "name": "r2",
  "version": "0.1.0",
  "optivor_api": "v1"
}
```

Fields:
- `name`: Unique provider driver identifier (lowercase, e.g. `r2`, `b2`, `gcs`).
- `version`: Driver binary version string.
- `optivor_api`: Target Optivor API contract version (currently `v1`).

### 2. Installing Drivers

Register a custom driver binary using the CLI:

```bash
optivor driver install /path/to/optivor-driver-r2
```

The CLI executes `--optivor-handshake`, validates the JSON response, and registers the driver in `~/.config/optivor/drivers.json`.

## Developer Guide Series

For a complete step-by-step guide to developing, testing, and contributing a custom provider driver:

1. [Architecture & Specification Overview](./developer-driver-overview.md)
2. [Step-by-Step Implementation Guide](./developer-driver-guide.md)
3. [Testing & Local Verification](./developer-driver-testing.md)
4. [Submission & Registry Contribution](./developer-driver-submission.md)

