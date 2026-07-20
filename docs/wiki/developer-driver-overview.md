# Developer Guide: Storage Driver Architecture & Specification

This document details the architectural specifications for building out-of-process storage drivers for Optivor.

## 1. Architectural Overview

Per **ADR-0010**, Optivor core maintains a provider-agnostic stance. The core runtime ships only with a universal S3-compatible driver. All provider-specific storage integrations (e.g., Cloudflare R2 native extensions, Backblaze B2, Google Cloud Storage, Azure Blob) are implemented as standalone external binaries following the `optivor-driver-<name>` naming convention.

```
+-------------------------------------------------------------+
|                      Optivor Core                           |
|  +--------------------+    +----------------------------+   |
|  |  Internal Storage  | -> | Out-of-Process Controller  |   |
|  +--------------------+    +----------------------------+   |
+------------------------------------------|------------------+
                                           | (exec / IPC)
                                           v
+-------------------------------------------------------------+
|             Standalone Provider Driver Executable           |
|                optivor-driver-<provider>                    |
+-------------------------------------------------------------+
```

### Key Benefits
- **Language Agnostic**: Drivers can be written in Go, Rust, Python, Node.js, C++, or any executable compiled language.
- **Isolation**: Crashes or resource spikes in third-party provider SDKs do not impair core Optivor runtime stability.
- **Independent Release Lifecycle**: Providers update independently of Optivor core versions.

---

## 2. Driver Contract & Handshake Protocol

Every driver executable MUST support the `--optivor-handshake` CLI flag.

### Handshake Invocation
```bash
optivor-driver-custom --optivor-handshake
```

### Handshake JSON Response Specification
The binary MUST write a valid JSON object to standard output (`stdout`) with an exit code of `0`:

```json
{
  "name": "r2",
  "version": "1.0.0",
  "optivor_api": "v1",
  "capabilities": ["read", "write", "stat", "delete", "presign"],
  "author": "Optivor Community",
  "homepage": "https://github.com/optivor/optivor-driver-r2"
}
```

#### Field Specifications:
| Field | Type | Required | Description |
|---|---|---|---|
| `name` | String | Yes | Lowercase unique provider identifier (alphanumeric and hyphens only). |
| `version` | String | Yes | Driver semantic version string (e.g., `1.0.0`). |
| `optivor_api` | String | Yes | Supported Optivor API protocol contract version (currently `v1`). |
| `capabilities` | Array[String] | Yes | List of supported operations: `read`, `write`, `stat`, `delete`, `presign`. |
| `author` | String | No | Driver author or organization name. |
| `homepage` | String | No | Repository URL or documentation link. |

---

## 3. Storage Operation Protocol (`v1`)

During execution, Optivor invokes the driver binary passing commands and parameters via command-line arguments or standard input JSON.

### Execution Context & Environment Variables
Optivor passes runtime configuration and credential secrets to the driver via environment variables:

- `OPTIVOR_DRIVER_PROVIDER`: Name of the targeted provider driver.
- `OPTIVOR_STORAGE_ENDPOINT`: Target storage service endpoint URL.
- `OPTIVOR_STORAGE_BUCKET`: Target bucket name.
- `OPTIVOR_STORAGE_ACCESS_KEY`: Storage access key / ID.
- `OPTIVOR_STORAGE_SECRET_KEY`: Storage secret key.
- `OPTIVOR_STORAGE_REGION`: Target storage region (if applicable).

### Standard Execution Command Interface
```bash
optivor-driver-<provider> <command> <key> [flags]
```

#### Supported Commands:
1. `get <key>`: Stream object payload to `stdout`.
2. `put <key>`: Read object payload from `stdin` and write to storage.
3. `head <key>`: Retrieve metadata (Content-Type, Content-Length, ETag) formatted as JSON to `stdout`.
4. `delete <key>`: Remove object from storage.

---

## 4. Exit Codes & Error Format

- Exit code `0`: Operation successful.
- Exit code `1`: Operational error (object not found, permission denied). Error message written to `stderr`.
- Exit code `2`: Invalid command arguments or invalid driver configuration.
