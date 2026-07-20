# Multi-Cloud & Multi-Bucket Management

Optivor supports managing multiple object storage buckets across different cloud providers (AWS S3, Cloudflare R2, Backblaze B2, Google Cloud Storage) simultaneously from a single instance.

## Declarative Multi-Bucket Configuration (`optivor.yaml`)

```yaml
buckets:
  - name: "primary-images"
    provider: s3
    endpoint: "https://s3.us-east-1.amazonaws.com"
    bucket: "my-aws-bucket"
    region: "us-east-1"
    access_key_id: "AWS_ACCESS_KEY"
    secret_access_key: "AWS_SECRET_KEY"
    access: public

  - name: "secure-assets"
    provider: r2
    endpoint: "https://<account-id>.r2.cloudflarestorage.com"
    bucket: "my-r2-bucket"
    access_key_id: "R2_ACCESS_KEY"
    secret_access_key: "R2_SECRET_KEY"
    access: signed
    fallback: "primary-images"

  - name: "internal-only"
    provider: b2
    endpoint: "https://s3.us-west-000.backblazeb2.com"
    bucket: "my-b2-bucket"
    access: private
```

## Per-Bucket Access Control Policies

Optivor enforces three levels of access control per bucket alias:

1. **`public`**: Anyone can request assets without authentication signatures.
2. **`signed`**: Requires HMAC-SHA256 URL signatures (`?sig=...&expires=...`).
3. **`private`**: Blocked externally; returns `403 Forbidden` on all requests.

## Failover Mechanisms

Configure a `fallback` key on any bucket to specify a backup bucket alias. If the primary storage provider returns a service error or goes down, Optivor automatically fetches the asset from the fallback bucket.
