# Granular IAM & Path-Level Access Control Guide

Optivor provides enterprise-grade Role-Based Access Control (RBAC) and path-level prefix authorization policies, enabling fine-grained security and path isolation for multi-tenant and multi-bucket deployments.

---

## Key Features

1. **Role-Based Privilege Definition (RBAC)**: Define declarative IAM roles (`admin`, `editor`, `reader-path-only`, or custom roles) with capability scopes (`read`, `write`, `lifecycle`, `*`).
2. **Path Prefix Authorization Policies**: Restrict API key tokens and identities to specific object key prefixes (e.g. `users/user-A/*`, `media/products/*`, `shared/logo.png`).
3. **Multi-Bucket Tenant Isolation**: Isolate user folders inside single or shared storage buckets without needing separate S3 infrastructure per tenant.
4. **Built-in & Custom Role Evaluation**: Seamlessly reference built-in roles (`admin`, `editor`, `reader-path-only`) or define custom roles in `optivor.yaml`.

---

## Configuration Reference (`optivor.yaml`)

```yaml
auth:
  # Declarative Role Definitions
  roles:
    - name: "admin"
      description: "Full system administration capabilities"
      capabilities: ["*"]
      allowed_paths: ["*"]
    - name: "editor"
      description: "Read and write access to media and asset folders"
      capabilities: ["read", "write"]
      allowed_paths: ["media/*", "assets/*"]
    - name: "tenant-user-a"
      description: "Restricted access to User A folder"
      capabilities: ["read"]
      allowed_paths: ["users/user-a/*"]

  # API Keys with IAM Role & Path Bindings
  api_keys:
    - key: "optivor_admin_secret_key"
      name: "System Admin"
      role: "admin"
      buckets: ["*"]

    - key: "optivor_user_a_key"
      name: "Tenant User A"
      role: "tenant-user-a"
      buckets: ["tenant-bucket"]
      allowed_paths: ["users/user-a/*"]

    - key: "optivor_marketing_key"
      name: "Marketing Team"
      role: "editor"
      buckets: ["primary-s3"]
      allowed_paths: ["marketing/*", "campaigns/*"]
```

---

## Path Matching Rules

Path matching rules evaluate requested object keys against configured `allowed_paths`:

| Pattern Format | Description | Example Target Key | Evaluation Result |
| :--- | :--- | :--- | :--- |
| `*` or empty | Allows access to all object paths | `any/folder/file.jpg` | **ALLOWED** |
| `prefix/*` | Allows any object key starting with `prefix/` | `users/user-a/avatar.png` | **ALLOWED** |
| `prefix/*` | Mismatched object prefix | `users/user-b/avatar.png` | **DENIED (403)** |
| `folder/file.ext` | Exact match rule | `shared/logo.png` | **ALLOWED** |

---

## HTTP Request Authentication Example

Supply the API key via `X-API-Key` or `Authorization: Bearer <key>` headers:

```bash
# Authorized request for User A path
curl -H "X-API-Key: optivor_user_a_key" \
  "https://optivor.example.com/image/tenant-bucket/users/user-a/avatar.png?w=200"

# Unauthorized request (attempting to access User B path) -> 403 Forbidden
curl -H "X-API-Key: optivor_user_a_key" \
  "https://optivor.example.com/image/tenant-bucket/users/user-b/avatar.png?w=200"
```

---

## Architecture Specification

For technical implementation details, middleware mechanics, and security guarantees, consult:
- [`docs/adr/0019-iam-path-privilege-and-role-policy.md`](../adr/0019-iam-path-privilege-and-role-policy.md)
