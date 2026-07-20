# Storage Driver Development & Distribution Guide

Optivor enforces a provider-agnostic core architecture (ADR-0002, ADR-0010). Storage drivers for cloud providers live outside the core repository as standalone binaries following the `optivor-driver-<name>` naming convention.

## Driver Specification & Handshake Protocol

An Optivor storage driver is an executable binary that satisfies the out-of-process protocol interface. When invoked with `--optivor-handshake`, the binary MUST print a JSON object to stdout (fd 1) and exit with code 0:

```json
{
  "name": "r2",
  "version": "1.0.0",
  "optivor_api": "v1"
}
```

Fields:
- `name`: Unique provider driver identifier (lowercase, e.g. `r2`, `b2`, `gcs`).
- `version`: Driver binary version string.
- `optivor_api`: Target Optivor API contract version (currently `v1`).

---

## Driver Installation & URL Formats

Optivor supports multiple installation sources via `optivor driver install`:

### 1. Local Binary File
Specify an absolute or relative path to a local driver binary:
```bash
optivor driver install /usr/local/bin/optivor-driver-r2
```

### 2. GitHub Shorthand Format
Specify a GitHub repository path with an optional `@tag` version:
```bash
# Installs latest release binary from GitHub repository
optivor driver install github:optivor/optivor-driver-r2

# Installs a specific tagged release version
optivor driver install github:optivor/optivor-driver-r2@v1.2.0
```

### 3. Direct HTTPS Release Binary URL
Specify a direct HTTPS URL to download a compiled driver binary:
```bash
optivor driver install https://github.com/optivor/optivor-driver-r2/releases/download/v1.2.0/optivor-driver-r2-linux-amd64
```

During installation, Optivor automatically verifies the binary, runs `--optivor-handshake`, validates the JSON output, and registers the metadata in `~/.config/optivor/drivers.json`.

---

## Release & Distribution Events

Driver authors distribute precompiled binaries via GitHub Releases:

1. **Tag & Release**: Tag repository releases using semantic versioning (`vX.Y.Z`).
2. **Release Assets**: Attach compiled binaries for target architectures (e.g. `optivor-driver-<name>-linux-amd64`).
3. **Automated CI**: Use GoReleaser or GitHub Actions to compile and attach assets on push events to `main` or `tags/v*`.

---

## Developer Guide Series

For a complete step-by-step guide to developing, testing, and contributing a custom provider driver:

1. [Architecture & Specification Overview](./developer-driver-overview.md)
2. [Step-by-Step Implementation Guide](./developer-driver-guide.md)
3. [Testing & Local Verification](./developer-driver-testing.md)
4. [Submission & Registry Contribution](./developer-driver-submission.md)
5. [Storage Driver SDK Specification](./driver-sdk-specification.md)
