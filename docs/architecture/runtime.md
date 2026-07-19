# Runtime Layer Architecture

## Overview

The Runtime layer (`internal/server`, `internal/pipeline`, `internal/cache`, `internal/storage`) is responsible for accepting HTTP requests, checking cache, fetching objects from S3 storage, performing image transformations using libvips, and returning response streams.

## Layer Boundaries and Dependency Directions

Per **ADR-0002**, dependency directions are strictly enforced:

```
cmd/optivor  →  internal/server
internal/server  →  internal/pipeline  →  internal/storage
internal/server  →  internal/cache
internal/pipeline  →  internal/config
```

### Constraints:
- `internal/pipeline` is **cache-unaware**. It executes `fetch -> transform -> encode` without referencing the cache layer.
- `internal/storage` has zero knowledge of `internal/server` or `internal/pipeline`.
- `internal/server` manages the request lifecycle, query validation, cache lookup, pipeline invocation, and cache population.

## Key Components

1. **HTTP Server (`internal/server`)**: Uses `chi` router to handle `GET /healthz` and wildcard image routes `GET /image/*`.
2. **Transform Pipeline (`internal/pipeline`)**: Uses `govips` for center crop cover, background-padded contain, and direct fill resizing, plus WebP encoding.
3. **Storage Driver (`internal/storage/s3`)**: Uses `minio-go` to stream objects from AWS S3 or MinIO backends.
4. **Filesystem Cache (`internal/cache/fs`)**: Uses SHA256 hashes of storage keys + transformation parameters to store cached image buffers atomically.
