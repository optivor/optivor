# CI/CD and Build Matrix Documentation

## Target Platforms (V0)

Optivor V0 targets **Linux (amd64 and arm64)** binary distributions.

### Cross-Compilation and cgo (libvips) Rationale

Optivor uses `davidbyttow/govips` (libvips C-bindings) for high-performance image resizing and WebP encoding.
Because libvips is a native C library dependency (`CGO_ENABLED=1`), cross-compiling for macOS (`darwin`) and Windows (`windows`) requires dedicated cross-compilation toolchains.

To maintain build reproducibility and deliver a reliable V0 release, non-Linux binary targets are deferred to **V0.1+**.

## Continuous Integration Workflow

GitHub Actions runs on every push and pull request targeting `staging` and `main`:
1. `go build ./...`
2. `go vet ./...`
3. `go test ./... -race -cover`

## Release Pipeline

GoReleaser cross-compiles native Linux binaries (`linux/amd64` and `linux/arm64`) and packages release archives and checksums automatically on release tags.
