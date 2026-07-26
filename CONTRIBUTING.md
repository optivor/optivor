# Contributing to Optivor

Thank you for your interest in contributing to **Optivor**! Optivor is an open-source image optimization engine and multi-ecosystem SDK suite designed for high-performance, secure media delivery.

This guide outlines our development workflow, repository organization, and pull request expectations.

---

## 1. Prerequisites & Local Environment

### Core Engine Requirements
- **Go**: Version `1.25.0` or higher
- **C Libraries (Image Codecs)**:
  - **Ubuntu / Debian**: `sudo apt-get install -y libvips-dev libheif-dev`
  - **macOS**: `brew install vips libheif`

### SDK Requirements (Monorepo)
- **Node.js**: `18+` (npm or pnpm) for JavaScript, React, Vue, and Next.js packages
- **Python**: `3.9+` (with `build` and `twine`) for the Python SDK
- **PHP**: `7.4+` (with Composer) for the PHP SDK

---

## 2. Project Architecture & Monorepo Structure

Optivor is structured as a modular Go core engine alongside client SDK packages:

```
optivor/
├── cmd/optivor/               # Binary CLI entrypoint
├── internal/
│   ├── server/                # HTTP transformation API & preset routing
│   ├── pipeline/              # Image resizing, cropping, & visual filters
│   ├── storage/               # Multi-bucket drivers (S3, GCS, Azure, Local)
│   └── cache/                 # Local, Redis, and Edge CDN cache adapters
├── packages/
│   ├── js/                    # @optivor/js (TypeScript / JS client)
│   ├── react/                 # @optivor/react (<OptivorImage /> component)
│   ├── vue/                   # @optivor/vue (Vue 3 / Nuxt component)
│   ├── next-loader/           # @optivor/next (Next.js custom loader)
│   ├── python/                # optivor (PyPI Python SDK)
│   └── php/                   # optivor/optivor-php (Composer PHP SDK)
└── docs/                      # Architectural Decision Records (ADRs) & Wiki
```

Before proposing architectural changes, please read the relevant ADRs:
- [`docs/adr/0000-project-scope.md`](./docs/adr/0000-project-scope.md) — Scope boundary
- [`docs/adr/0001-project-philosophy.md`](./docs/adr/0001-project-philosophy.md) — Design principles
- [`docs/adr/0002-high-level-architecture.md`](./docs/adr/0002-high-level-architecture.md) — Subsystem boundaries

---

## 3. Branching & Commit Conventions

### Branching Strategy
Optivor follows a strict **dual-branching model**:
- **`main`**: Production release branch. Contains tagged releases and stable code.
- **`staging`**: Integration branch. **All Pull Requests must target `staging`**.

When creating a feature or bug fix:
```bash
git checkout staging
git pull origin staging
git checkout -b feat/your-feature-name
```

### Commit Message Format
We follow Conventional Commits:
- `feat(scope)`: A new feature (e.g. `feat(server): add focal point crop parameter`)
- `fix(scope)`: A bug fix (e.g. `fix(storage/s3): resolve bucket path escaping`)
- `docs(scope)`: Documentation changes (e.g. `docs(readme): update preset usage`)
- `bump(packages)`: Version bump across SDKs
- `refactor(scope)`: Code refactoring without behavioral changes

---

## 4. Local Building & Testing

Run local tests and security scanners before submitting a Pull Request:

```bash
# Build binary
make build

# Run Go unit & integration tests
go test ./... -v -race -cover

# Run static analysis & vulnerability scans
go vet ./...
golangci-lint run
govulncheck ./...

# Test JavaScript SDKs
cd packages/js && npm test
```

---

## 5. Subsystem Map & Guidance

| I want to... | Start Here | Notes |
| :--- | :--- | :--- |
| **Fix a Core Bug** | Open a **Bug Report** issue first | Unless trivial, discuss bug details before opening a PR |
| **Add Storage Driver** | Read [`docs/adr/0003-extensibility-model.md`](./docs/adr/0003-extensibility-model.md) | Submit a **Driver Proposal** issue before writing code |
| **Propose Engine Feature** | Open a **Feature Request** issue | Ensure fit with current milestone scope in `ROADMAP.md` |
| **Improve SDKs / Docs** | PRs directly to `staging` | Include updated tests and `README.md` examples |

---

## 6. Code of Conduct

Optivor is committed to providing a welcoming, inclusive, and harassment-free community. Please review and adhere to our [`CODE_OF_CONDUCT.md`](./CODE_OF_CONDUCT.md).
