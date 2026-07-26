---
name: "🐛 Bug Report"
about: Report an issue in the Optivor Go engine, server runtime, CLI, or client SDKs
title: "[Bug]: "
labels: bug
assignees: ''
---

## 🐛 Description
A clear and concise description of what the bug is.

## 📦 Affected Subsystem / Component
Please select the component where the issue occurs:
- [ ] Go Server Engine (`internal/server`, `internal/pipeline`)
- [ ] Storage Driver (`S3`, `GCS`, `Azure`, `Local`)
- [ ] Cache Adapter (`Local`, `Redis`, `Edge CDN`)
- [ ] CLI Tool (`cmd/optivor`, `internal/cli`)
- [ ] `@optivor/js` (JS / TS SDK)
- [ ] `@optivor/react` (React SDK)
- [ ] `@optivor/vue` (Vue SDK)
- [ ] `@optivor/next` (Next.js Loader)
- [ ] `optivor` (Python SDK)
- [ ] `optivor/optivor-php` (PHP SDK)

## 🔄 Steps to Reproduce
1. Go to '...'
2. Request URL / Run command '...'
3. See error '...'

## 🎯 Expected Behavior
A clear description of what you expected to happen instead.

## 💻 Environment & Version
- **Optivor Engine Version / Commit**: (e.g. `v1.2.2` or git commit `abcdef`)
- **SDK Package & Version**: (e.g. `@optivor/react@1.2.2`)
- **Storage Backend**: (e.g., AWS S3, MinIO, Google Cloud Storage, Local Filesystem)
- **Deployment Platform**: (Docker, Kubernetes, Standalone Binary, Cloudflare Workers)
- **OS & Architecture**: (e.g. Ubuntu 24.04 ARM64, macOS Sequoia M2, Windows 11)

## ⚙️ Configuration (Redact any API Keys / Secrets)
```yaml
# Relevant optivor.yaml snippet
```

## 📜 Logs & Error Output
```text
Paste relevant server logs or console stack trace here
```

## 🖼️ Screenshots / Minimal Code Reproduction
If applicable, add screenshots or a minimal code snippet to help explain your problem.
