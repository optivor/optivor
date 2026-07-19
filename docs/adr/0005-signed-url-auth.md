# ADR-0005: Signed URL & Authorization

## Status

Accepted

## Context

V0 released Optivor as an unauthenticated HTTP service. To expose Optivor safely to the public internet without relying on an external reverse-proxy auth layer, Optivor needs built-in authorization for image requests.

Per ADR-0002, the runtime layer (`internal/server`) handles HTTP routing and request processing, while the image processing pipeline (`internal/pipeline`) must remain completely cache-unaware and auth-unaware. Per ADR-0003, authentication and policy checks operate in the runtime hot path as middleware.

## Decision

1. **Signature Scheme**: Optivor adopts HMAC-SHA256 signature verification via HTTP query parameters:
   - `expires`: Unix timestamp (in seconds) after which the request is invalid.
   - `sig`: Hex-encoded HMAC-SHA256 hash of `path + "?expires=" + expires` signed using `auth.signed_urls.secret`.

2. **Configuration**:
   - `auth.signed_urls.enabled` (boolean, default: `false` for backward compatibility with V0).
   - `auth.signed_urls.secret` (string, loaded from `OPTIVOR_AUTH_SECRET` environment variable or YAML configuration).
   - `auth.signed_urls.max_age` (integer seconds, default: 3600).

3. **HTTP Behavior**:
   - When `enabled: false`: all requests pass unauthenticated (V0 behavior).
   - Valid signature and `expires >= now`: returns HTTP `200 OK`.
   - Expired signature (`expires < now`): returns HTTP `403 Forbidden`.
   - Invalid or missing signature: returns HTTP `401 Unauthorized`.

4. **Layer Boundary Isolation**:
   - Signed URL validation is implemented strictly inside `internal/server` as HTTP middleware.
   - `internal/pipeline`, `internal/storage`, and `internal/cache` remain completely unaware of auth logic.

## Consequences

- Direct public internet exposure of Optivor becomes safe when `auth.signed_urls.enabled` is enabled.
- Clients must compute valid HMAC-SHA256 signatures when issuing requests if signed URLs are enabled.
