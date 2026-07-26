# ADR-0019: IAM Path Privilege and Role Policy Architecture

- **Status**: Accepted
- **Deciders**: Optivor Core Architecture Team
- **Date**: 2026-07-26

---

## Context and Problem Statement

As Optivor scales in enterprise multi-tenant and B2B deployment scenarios, single bucket authorization or static API keys are insufficient to meet granular access control requirements:
1. **Path-Based Prefix Authorization**: Tenants and users require path-level isolation within shared S3/R2 storage buckets (e.g., `user-A` limited strictly to `s3://bucket/users/user-A/*`, preventing unauthorized reads or modifications of `s3://bucket/users/user-B/*`).
2. **Role-Based Access Control (RBAC)**: Administrators need declarative role definitions (`admin`, `editor`, `reader-path-only`) that map identities to specific operation capability scopes (`read`, `write`, `lifecycle`, `*`) and object key prefixes.
3. **Signed URL Delegation**: Signed URL signatures and delegation tokens must support embedding or evaluating path prefix constraints to enforce tamper-proof folder-level access policies.

---

## Decision Outcome

We decide to implement Granular IAM & Path-Level Access Control across Optivor's configuration, authentication middleware, and authorization engine:

### 1. Declarative IAM Role Definitions (`auth.roles`)
Define reusable IAM roles in `optivor.yaml`:
```yaml
auth:
  roles:
    - name: "admin"
      description: "Full system administration capabilities"
      capabilities: ["*"]
      allowed_paths: ["*"]
    - name: "editor"
      description: "Read and write access to media directories"
      capabilities: ["read", "write"]
      allowed_paths: ["media/*", "uploads/*"]
    - name: "reader-path-only"
      description: "Read-only access restricted to specific user folders"
      capabilities: ["read"]
      allowed_paths: ["user-A/*"]
```

### 2. Enhanced API Key Bindings (`auth.api_keys`)
API keys can directly bind to declared roles or specify explicit allowed path prefixes and scopes:
```yaml
auth:
  api_keys:
    - key: "optivor_secret_key_123"
      name: "User-A Access Key"
      role: "reader-path-only"
      buckets: ["primary-s3"]
      allowed_paths: ["user-A/*"]
```

### 3. Path Prefix & Capabilities Matching Rules
- **Path Prefix Validation**:
  - `*` or empty matches any key path.
  - `folder/*` or `prefix/*` matches any object key starting with `folder/` or `prefix/`.
  - Exact path string matches only that exact key.
- **Capabilities Matching**:
  - `*` grants all operations (`read`, `write`, `lifecycle`).
  - Specific capability strings gate operation types accordingly.

### 4. Integration with Authentication Middleware
`ValidateAPIKey` and `ValidateIAMAccess` evaluate both direct key properties and inherited role rules:
1. Match API Key from `X-API-Key` or `Authorization: Bearer <key>` header.
2. Resolve referenced `Role` (if specified in key config).
3. Validate requested bucket against allowed bucket list (`*` or explicit bucket name).
4. Validate requested capability scope (`read`, `write`, etc.).
5. Validate target object key against merged allowed path prefixes (`AllowedPaths`).

---

## Consequences

### Positive
- Enterprise-grade tenant isolation and RBAC without requiring separate S3 buckets per user.
- Declarative configuration in `optivor.yaml` with automatic environment variable overrides.
- Seamless compatibility with dynamic signed URL delegation and multi-bucket routing (`internal/storage/router`).

### Negative
- Mild computational overhead for path pattern matching on authenticated requests, optimized via fast prefix comparison strings.
