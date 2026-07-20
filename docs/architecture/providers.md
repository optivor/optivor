# Architecture: Storage Providers & Ecosystem

## Core Principles

Per **ADR-0002** (High-Level Architecture) and **ADR-0010** (Provider Registry & Driver Convention):

1. **Vendor Independence**: Core Optivor contains zero vendor-specific SDKs or cloud platform dependencies.
2. **Universal S3 Baseline**: S3 protocol compliance is used as the universal baseline driver (`internal/storage/s3`). AWS S3, MinIO, Cloudflare R2, Backblaze B2, and Wasabi are all supported via S3 API compatibility.
3. **Out-of-Process Custom Drivers**: Custom drivers with specialized cloud SDK interactions ship as standalone binaries (`optivor-driver-<name>`) managed via `optivor driver install`.

## Dependency Direction

```
Runtime → Storage Driver Interface (in-process S3 or out-of-process subprocess)
```

The storage layer NEVER imports runtime packages, cache logic, or CLI code.
