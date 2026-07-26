# Security Policy

Optivor is built for mission-critical media delivery. We maintain a zero-trust approach to image URL parsing, storage driver access, and transformation pipeline security.

---

## Supported Versions

Optivor actively supports security updates for the following release tracks:

| Version Line | Supported | Description |
| :--- | :---: | :--- |
| **`v1.2.x`** | :white_check_mark: | Active production release (Preset routing, Focal Crop, Overlays) |
| **`v1.1.x`** | :white_check_mark: | Security patches & critical vulnerability backports |
| **`v1.0.x`** | :white_check_mark: | Security patches & critical bug fixes |
| **`< v1.0`** | :x: | End of Life (EOL) - Please upgrade to `v1.2.x` |

---

## Reporting a Vulnerability

We take the security of Optivor seriously. If you discover a security vulnerability, please **DO NOT open a public issue**.

### Reporting Channel

Please submit security reports through **GitHub Private Vulnerability Reporting**:
1. Go to the Optivor GitHub repository: [github.com/optivor/optivor](https://github.com/optivor/optivor)
2. Click on the **Security** tab -> **Report a vulnerability**.
3. If Private Vulnerability Reporting is unavailable, contact the core security maintainers listed in [`CODEOWNERS`](./CODEOWNERS).

### Details to Include in Your Report

To help us triage and patch the issue quickly, please include:
- **Subsystem / Component**: (e.g. Go Engine, `internal/storage/local`, S3 driver, JS/React/Vue/Python/PHP SDK, CLI).
- **Vulnerability Type**: (e.g. SSRF, Path Traversal, Denial of Service, Signature Bypass, Memory Leak).
- **Proof of Concept (PoC)**: Minimal step-by-step reproduction steps or HTTP request curl commands.
- **Impact Assessment**: Estimated severity and potential threat vector.

---

## Vulnerability Handling & Response SLA

- **Initial Acknowledgment**: Within **24 to 48 hours**.
- **Triage & Severity Assessment**: Within **3 business days**.
- **Private Patch Development**: Prepared on a private security release branch.
- **CVE & Advisory Disclosure**: Published via GitHub Security Advisories alongside a patch version release.

---

## Core Security Architectural Guarantees

Optivor enforces defense-in-depth across the engine:

1. **HMAC URL Signatures**: When `security_key` is enabled, all incoming transform parameters and paths are validated using HMAC-SHA256 tokens (`s=...`). Any URL parameter tampering or parameter stripping returns an instant `403 Forbidden`.
2. **Path Traversal Protection (CWE-22)**: Storage drivers validate and resolve relative keys against root directory basenames to prevent path escape (`../`).
3. **SSRF Mitigation**: Remote fetch adapters validate destination IP ranges against private CIDR blocks (`10.0.0.0/8`, `192.168.0.0/16`, `172.16.0.0/12`, `127.0.0.1/8`) to prevent Server-Side Request Forgery.
4. **Watermark Anti-Tamper Enforcement**: Watermark overlays defined via server-side `/preset/{name}/{key}` routes reside strictly on the server, preventing clients from removing overlay parameters.
