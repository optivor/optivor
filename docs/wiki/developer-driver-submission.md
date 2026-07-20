# Developer Guide: Driver Submission & Registry Contribution

This document outlines the step-by-step process for submitting your completed custom storage driver to the official Optivor Provider Registry.

---

## Step 1: Submission Requirements & Guidelines

Before submitting your driver for inclusion, ensure it fulfills the following guidelines:

1. **Repository Naming**: Public GitHub repository named `optivor-driver-<provider>` (e.g., `optivor-driver-r2`).
2. **License**: Open-source license (MIT, Apache 2.0, or BSD-3-Clause).
3. **Cross-Platform Release**: Pre-compiled binary releases for `linux/amd64`, `linux/arm64`, `darwin/amd64`, `darwin/arm64`, and `windows/amd64` uploaded to GitHub Releases.
4. **Documentation**: Clear `README.md` explaining provider prerequisites, configuration variables, and version support.
5. **Handshake Compliant**: Implements `--optivor-handshake` returning protocol version `v1`.

---

## Step 2: Continuous Integration & Release Packaging

Configure GitHub Actions in your driver repository to build cross-compiled binaries on every tagged release.

### Sample `.github/workflows/release.yml`
```yaml
name: Release Driver Binaries

on:
  push:
    tags:
      - 'v*'

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'
      - name: Build cross-platform binaries
        run: |
          GOOS=linux GOARCH=amd64 go build -o bin/optivor-driver-b2-linux-amd64 ./cmd/optivor-driver-b2
          GOOS=linux GOARCH=arm64 go build -o bin/optivor-driver-b2-linux-arm64 ./cmd/optivor-driver-b2
          GOOS=darwin GOARCH=arm64 go build -o bin/optivor-driver-b2-darwin-arm64 ./cmd/optivor-driver-b2
      - name: Upload Binaries to Release
        uses: softprops/action-gh-release@v1
        with:
          files: bin/*
```

---

## Step 3: Submitting to the Registry Index

To list your driver in the official registry (so users can discover and install it via `optivor driver install <provider>`), submit a Pull Request to the core Optivor repository or registry index.

### 1. Open a PR to `staging`
Branch off `staging` in `optivor/optivor` and update `docs/wiki/storage-drivers.md` (or registry catalog):

```markdown
| Provider | Driver Name | Repository | Maintainer | Protocol |
|---|---|---|---|---|
| Backblaze B2 | `b2` | [optivor-driver-b2](https://github.com/example/optivor-driver-b2) | Community | v1 |
```

### 2. PR Information Checklist
Include the following in your Pull Request description:
- [x] Driver Repository Link
- [x] Tested Provider Services / Endpoints
- [x] Handshake Output snippet
- [x] Verification test results

---

## Step 4: Review & Automated Validation

Once your PR is submitted:
1. **Automated Handshake Check**: CI runs `optivor driver install` against the release binary to verify handshake output.
2. **CODEOWNERS Review**: An Optivor maintainer reviews licensing, security, and protocol compliance.
3. **Registry Listing**: Upon approval and merge into `staging`, your driver will be listed as an official community provider!
