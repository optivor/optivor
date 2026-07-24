# Configuration Reference

Optivor uses Viper for configuration parsing from `optivor.yaml` and environment variables.

## Full `optivor.yaml` Schema

```yaml
server:
  port: 8080                  # HTTP server port
  read_timeout: 30s           # Max duration for reading request
  write_timeout: 30s          # Max duration for writing response
  request_timeout: 30s        # Total request timeout context
  log_level: "info"           # debug | info | warn | error
  log_format: "text"          # text | json
  rate_limit:
    enabled: true             # Enable per-IP token bucket rate limiter
    rps: 10                   # Requests per second per IP
    burst: 20                 # Token bucket burst capacity
  image:
    max_width: 5000           # Maximum allowed transform width (px)
    max_height: 5000          # Maximum allowed transform height (px)

remote:
  enabled: true               # Enable transparent remote URL fetching (/fetch, /remote)
  allowed_domains:            # Domain whitelist for remote fetching ("*" for all)
    - "example.com"
    - "cdn.mydomain.com"

presets:
  avatar:                     # Preset name (/preset/avatar/*)
    w: 150
    h: 150
    f: "webp"
    fit: "cover"

crawler:
  enabled: true               # Enable bot & crawler detection/throttling
  max_concurrency_per_variant: 5 # Concurrency limit per image variant transformation

storage:
  driver: "s3"                # Storage driver selection (s3, minio, r2)
  s3:
    endpoint: "https://..."   # S3 endpoint URL
    bucket: "my-bucket"       # S3 bucket name
    region: "us-east-1"       # S3 region
    access_key_id: "..."      # S3 Access Key ID
    secret_access_key: "..."  # S3 Secret Access Key

buckets:
  - name: "primary-images"
    provider: s3
    endpoint: "https://s3.us-east-1.amazonaws.com"
    bucket: "my-aws-bucket"
    region: "us-east-1"
    access: public

  - name: "secure-assets"
    provider: r2
    endpoint: "https://account-id.r2.cloudflarestorage.com"
    bucket: "my-r2-bucket"
    access: signed
    fallback: "primary-images"

cache:
  fs:
    dir: "/tmp/optivor-cache" # Cache directory path
    max_size_mb: 1024         # LRU eviction limit in megabytes

telemetry:
  enabled: true               # Enable OpenTelemetry tracing
  otlp_endpoint: ""           # OTLP/gRPC collector target URL
  service_name: "optivor"     # Service identifier in traces
  sampling_ratio: 1.0         # Trace sampling probability (0.0 to 1.0)

auth:
  signed_urls:
    enabled: false            # Require HMAC signatures on request URLs
    secret: ""                # HMAC signing secret key
    max_age: 3600             # Default URL expiration age in seconds

image:
  contain_background_color: "#ffffff"  # Background color for fit=contain
  max_pixels: 25000000                 # Decompression-bomb pixel threshold (~5000x5000)
  max_decode_mb: 64                    # libvips startup memory ceiling
```

## Environment Overrides

Any configuration field can be overridden using `OPTIVOR_` environment variables with underscore delimiters:

- `OPTIVOR_STORAGE_S3_SECRET_ACCESS_KEY` → `storage.s3.secret_access_key`
- `OPTIVOR_AUTH_SIGNED_URLS_SECRET` → `auth.signed_urls.secret`
- `OPTIVOR_SERVER_PORT` → `server.port`
- `OPTIVOR_REMOTE_ENABLED` → `remote.enabled`
- `OPTIVOR_CRAWLER_ENABLED` → `crawler.enabled`
