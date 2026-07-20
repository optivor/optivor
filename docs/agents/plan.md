# Optivor — Uygulama Planı

## Genel Bakış

Bu belge Optivor'un milestone bazlı uygulama planıdır.
Her milestone kendi bölümünde tanımlanır; adımlar AGENTS.md §1–§5 kurallarını
birebir izler (branch → commit → PR → CI yeşil → merge sırası).

**Tamamlanma kriterleri ve teknoloji seçimleri ADR'lerden devralınır; bu belgede
tekrarlanmaz — yalnızca "ne yapılacak ve hangi sırada" bilgisi tutulur.**

---

## ✅ V0 — Prove the Core Loop (TAMAMLANDI)

**Hedef:** S3-uyumlu bir bucket'tan görüntü alıp WebP'ye dönüştüren, tek standalone
Go binary'si üretmek.

**Durum:** `staging` branch'inde tüm PR'lar merge edildi. `main`'e release PR
açılmamışsa açılması gerekir (bkz. §"Release" adımı).

### Tamamlanan PR'lar (staging'de birleşti)

| PR | Branch | İçerik |
|---|---|---|
| #1 | `feat/v0-skeleton-and-config` | Go modülü, dizin iskeleti, Viper config |
| #2 | `feat/storage-s3-driver` | `internal/storage/s3` — minio-go driver |
| #3 | `feat/pipeline-govips` | `internal/pipeline` — govips transform |
| #4 | `feat/server-and-cache` | `internal/server` + `internal/cache/fs` |
| #5 | `feat/cmd-main-and-e2e` | `cmd/optivor`, E2E DoD testi, Dockerfile, CI |

### V0 Release Adımı (henüz yapılmadıysa)

```bash
# staging yeşilse release PR aç:
git checkout staging && git pull origin staging
gh pr create \
  --base main \
  --head staging \
  --title "release: promote staging → main (v0.1.0-alpha)" \
  --body "V0 milestone complete. See commit log for full scope."

# PR merge edildikten sonra tag:
git checkout main && git pull origin main
git tag v0.1.0-alpha
git push origin v0.1.0-alpha   # GoReleaser tetikler
```

---

## ✅ V0.1 — Make it Safe to Expose (TAMAMLANDI)

**ROADMAP.md hedefleri:**
- Signed URL ve policy-based authorization
- AVIF output (govips zaten destekler — encoder dependency kararı)
- Resource limits: max pixel count, max memory per transform, request timeouts (decompression-bomb koruması)
- Basic rate limiting

**V0'dan devralınan borçlar (V0.1'de kapatılacak):**
- Cache eviction / `max_size_mb` (S4 — V0.1'e ertelendi)
- macOS/Windows binary (S5 — cgo cross-compile çözümü)

---

### V0.1 Adım Sırası (AGENTS.md §5 kuralı)

> Her adım = bir `feat/...` branch + bir PR → `staging`.
> Bağımsız adımlar paralel dallanabilir; bağımlı adımlar sıralı açılır.

---

#### Adım A — ADR-0005: Signed URL & Authorization

**Neden ADR önce?** Auth policy yeni bir extension point — ADR-0003'ün "yeni
extension point = önce ADR" kuralı.

```bash
git checkout staging && git pull origin staging
git checkout -b docs/adr-0005-signed-url-auth
```

**Commit:**
```
docs(adr): add ADR-0005 for signed URL and policy-based authorization
```

**ADR-0005 kapsamı (minimal):**
- URL imzalama algoritması: HMAC-SHA256, query param `?sig=&expires=`
- Secret key kaynağı: `OPTIVOR_AUTH_SECRET` env var
- Policy: sadece zaman bazlı expiry (V0.1); scope/path kısıtlaması V0.3+
- Middleware yerleşimi: `internal/server` — pipeline habersiz kalır

**PR:**
```bash
git push -u origin docs/adr-0005-signed-url-auth
gh pr create --base staging \
  --title "docs(adr): add ADR-0005 for signed URL authorization" \
  --body-file .github/PULL_REQUEST_TEMPLATE.md
```

---

#### Adım B — Signed URL implementasyonu

> **Ön koşul:** Adım A (ADR-0005) merge edilmiş olmalı.

```bash
git checkout staging && git pull origin staging
git checkout -b feat/signed-url-auth
```

**Commit sırası (her biri ayrı commit):**

```
feat(server): add HMAC-SHA256 signed URL middleware

feat(config): add auth.secret config field with OPTIVOR_AUTH_SECRET env

test(server): add table-driven tests for signed URL validation

docs(readme): document signed URL usage and security warning removal
```

**Config değişikliği (`optivor.yaml`):**
```yaml
auth:
  signed_urls:
    enabled: false          # V0 geriye dönük uyumlu — default off
    secret: ""              # env: OPTIVOR_AUTH_SECRET
    max_age: 3600           # saniye cinsinden default expiry
```

**HTTP davranışı:**
| Durum | Status |
|---|---|
| `signed_urls.enabled: false` | Tüm istekler geçer (V0 davranışı) |
| Geçerli imza | `200 OK` |
| Süresi dolmuş imza | `403 Forbidden` |
| Geçersiz/eksik imza | `401 Unauthorized` |

**PR:**
```bash
git push -u origin feat/signed-url-auth
gh pr create --base staging \
  --title "feat(server): implement signed URL authorization middleware" \
  --body-file .github/PULL_REQUEST_TEMPLATE.md
```

---

#### Adım C — Resource Limits & Decompression-Bomb Koruması

> **Ön koşul:** Yok — Adım B ile paralel açılabilir.

```bash
git checkout staging && git pull origin staging
git checkout -b feat/resource-limits
```

**Commit sırası:**

```
feat(config): add image.max_pixels and server.request_timeout config fields

feat(pipeline): enforce max pixel count before transform (decompression-bomb guard)

feat(server): add per-request timeout via context deadline

test(pipeline): add tests for oversized image rejection

test(server): add timeout integration test
```

**Config değişikliği:**
```yaml
server:
  request_timeout: 30s     # toplam istek timeout'u

image:
  max_pixels: 25000000     # ~5000×5000px — pipeline girişinde kontrol
  max_decode_mb: 64        # libvips memory ceiling (govips.Startup options)
```

**HTTP hata tablosu (yeni satırlar):**
| Durum | Status |
|---|---|
| Kaynak görüntü `max_pixels` aşıyor | `413 Content Too Large` |
| İstek `request_timeout` aşıyor | `408 Request Timeout` |
| libvips memory limit aşılırsa | `500 Internal Server Error` (loglanır) |

**PR:**
```bash
git push -u origin feat/resource-limits
gh pr create --base staging \
  --title "feat(pipeline,server): add resource limits and decompression-bomb protection" \
  --body-file .github/PULL_REQUEST_TEMPLATE.md
```

---

#### Adım D — Basic Rate Limiting

> **Ön koşul:** Yok — Adım B ve C ile paralel açılabilir.

```bash
git checkout staging && git pull origin staging
git checkout -b feat/rate-limiting
```

**Karar:** `golang.org/x/time/rate` — sıfır dış bağımlılık, token bucket,
ADR-0001 "minimal dependencies" ilkesiyle uyumlu. Dağıtık rate limit V0.3+.

**Commit sırası:**

```
feat(config): add server.rate_limit config section

feat(server): implement per-IP token bucket rate limiter middleware

test(server): add rate limit middleware tests with burst simulation

docs(readme): document rate limiting configuration
```

**Config değişikliği:**
```yaml
server:
  rate_limit:
    enabled: true
    rps: 10          # saniyede istek sayısı (per-IP)
    burst: 20        # burst kapasitesi
```

**HTTP davranışı:**
| Durum | Status |
|---|---|
| Limit aşılır | `429 Too Many Requests` + `Retry-After` header |

**PR:**
```bash
git push -u origin feat/rate-limiting
gh pr create --base staging \
  --title "feat(server): add per-IP token bucket rate limiting" \
  --body-file .github/PULL_REQUEST_TEMPLATE.md
```

---

#### Adım E — AVIF Output

> **Ön koşul:** Yok — diğer adımlarla paralel açılabilir.

```bash
git checkout staging && git pull origin staging
git checkout -b feat/avif-output
```

**Karar:** `govips` zaten AVIF encode destekliyor (libvips ≥ 8.13 + libheif).
Dockerfile'a `libheif-dev` eklenir. Cross-compile etkilenmez (aynı CGO_ENABLED=1
Linux-only matrix devam eder).

**Commit sırası:**

```
chore(docker): add libheif-dev to Dockerfile for AVIF support

feat(pipeline): add AVIF encode path in encode.go

feat(server): add 'avif' to format query param whitelist

test(pipeline): add AVIF encode unit test

docs(ci-cd): note libheif dependency addition
```

**Config değişikliği:** Yok — `?format=avif` query param olarak gelir.

**HTTP whitelist güncellemesi (`server.go`):**
```go
// V0: whitelist = {"webp", ""}
// V0.1: whitelist = {"webp", "avif", ""}
```

**PR:**
```bash
git push -u origin feat/avif-output
gh pr create --base staging \
  --title "feat(pipeline): add AVIF output support via govips/libheif" \
  --body-file .github/PULL_REQUEST_TEMPLATE.md
```

---

#### Adım F — Cache Eviction (V0 Borcu)

> **Ön koşul:** Yok — paralel açılabilir.

```bash
git checkout staging && git pull origin staging
git checkout -b feat/cache-eviction
```

**Karar:** LRU eviction, saf Go (`container/list`) — dış bağımlılık yok.
`max_size_mb` config alanı aktif edilir.

**Commit sırası:**

```
feat(cache/fs): implement LRU eviction with max_size_mb limit

feat(config): activate cache.fs.max_size_mb field (was deferred in V0)

test(cache/fs): add eviction tests for size threshold

docs(readme): remove V0 cache growth warning (eviction now implemented)
```

**Config değişikliği:**
```yaml
cache:
  fs:
    dir: "/tmp/optivor-cache"
    max_size_mb: 512    # V0.1'de aktif
```

**PR:**
```bash
git push -u origin feat/cache-eviction
gh pr create --base staging \
  --title "feat(cache/fs): implement LRU eviction with max_size_mb limit" \
  --body-file .github/PULL_REQUEST_TEMPLATE.md
```

---

#### Adım G — macOS/Windows Binary (cgo cross-compile)

> **Ön koşul:** Adım E (AVIF) merge edilmiş olmalı — Dockerfile değişikliğiyle çakışmamak için.

```bash
git checkout staging && git pull origin staging
git checkout -b chore/goreleaser-cross-compile
```

**Yöntem:** `zig cc` cross-compile toolchain (GoReleaser + zig) — macOS arm64/amd64,
Windows amd64. CI'da `amd64-linux`'ta çalışır; ayrı runner gerekmez.

**Commit sırası:**

```
chore(build): add zig-based cross-compile toolchain for macOS and Windows targets

chore(goreleaser): expand build matrix to darwin/amd64, darwin/arm64, windows/amd64

chore(ci): add cross-compile smoke test to GitHub Actions

docs(ci-cd): document cross-compile toolchain change and V0 deferral resolution
```

**PR:**
```bash
git push -u origin chore/goreleaser-cross-compile
gh pr create --base staging \
  --title "chore(build): add macOS and Windows cross-compile via zig toolchain" \
  --body-file .github/PULL_REQUEST_TEMPLATE.md
```

---

#### Adım H — V0.1 E2E Kabul Testi Güncellemesi

> **Ön koşul:** A–G'nin tamamı merge edilmiş olmalı.

```bash
git checkout staging && git pull origin staging
git checkout -b test/v01-e2e-acceptance
```

**Commit:**
```
test(e2e): extend DoD acceptance test for V0.1 features

Covers: signed URL validation, rate limit 429, AVIF output,
oversized image 413, and cache eviction under load.
```

**Test senaryoları:**
- `GET /image/...?format=avif` → `Content-Type: image/avif` ✓
- Geçersiz imzayla istek → `401 Unauthorized` ✓
- Burst istek → `429 Too Many Requests` ✓
- `max_pixels` aşan görüntü → `413` ✓
- Cache dolunca eviction — disk boyutu limit altında kalır ✓

**PR:**
```bash
git push -u origin test/v01-e2e-acceptance
gh pr create --base staging \
  --title "test(e2e): V0.1 acceptance test suite" \
  --body-file .github/PULL_REQUEST_TEMPLATE.md
```

---

#### V0.1 Release

```bash
git checkout staging && git pull origin staging
gh pr create \
  --base main \
  --head staging \
  --title "release: promote staging → main (v0.1.0)" \
  --body "V0.1 milestone: signed URLs, AVIF, rate limiting, resource limits, cache eviction, cross-compile."

# Merge sonrası:
git checkout main && git pull origin main
git tag v0.1.0
git push origin v0.1.0
```

---

### V0.1 Kapsam Dışı

| Özellik | Hedef |
|---|---|
| CLI (`optivor init`, `optivor deploy`) | V0.3 |
| Redis / distributed cache | Belirsiz |
| Scope/path bazlı auth policy | V0.3 |
| Prometheus metrics | V0.2 |
| Deployment adapter | V0.2 |

---

## 📋 V0.2 — First Deployment Adapter

**ROADMAP.md hedefleri:**
- Bir gerçek Deployment Adapter (hedef: standalone-as-systemd-service **veya** Fly.io)
- Structured logging ve temel metrics (Prometheus)

---

### V0.2 Adım Sırası

---

#### Adım I — Structured Logging (log/slog → JSON)

> **Ön koşul:** V0.1 tüm PR'ları merge edilmiş, `main`'de `v0.1.0` tagı var.

```bash
git checkout staging && git pull origin staging
git checkout -b feat/structured-logging
```

**Karar:** `log/slog` zaten mevcut. V0.2'de JSON handler eklenir; log level
config'den alınır. Dış bağımlılık yok.

**Commit sırası:**

```
feat(config): add server.log_level and server.log_format config fields

feat(server): switch slog handler to JSON when log_format=json

test(server): verify JSON log output format in integration test

docs(readme): document log_level and log_format configuration
```

**Config değişikliği:**
```yaml
server:
  log_level: "info"      # debug | info | warn | error
  log_format: "text"     # text | json
```

**PR:**
```bash
git push -u origin feat/structured-logging
gh pr create --base staging \
  --title "feat(server): add structured JSON logging with configurable level" \
  --body-file .github/PULL_REQUEST_TEMPLATE.md
```

---

#### Adım J — Prometheus Metrics

> **Ön koşul:** Adım I merge edilmiş olmalı.

```bash
git checkout staging && git pull origin staging
git checkout -b feat/prometheus-metrics
```

**Bağımlılık:** `github.com/prometheus/client_golang` — ADR-0004'te zaten seçilmiş.

**Commit sırası:**

```
chore(deps): add prometheus/client_golang dependency

feat(server): expose /metrics endpoint with Prometheus handler

feat(server): instrument request count, latency, and cache hit/miss metrics

test(server): add metrics endpoint smoke test

docs(readme): document /metrics endpoint and Prometheus scrape config
```

**Metrikler:**

| Metric | Type | Labels |
|---|---|---|
| `optivor_requests_total` | Counter | `status`, `format`, `fit` |
| `optivor_request_duration_seconds` | Histogram | `status` |
| `optivor_cache_hits_total` | Counter | — |
| `optivor_cache_misses_total` | Counter | — |
| `optivor_transform_duration_seconds` | Histogram | `format`, `fit` |

**HTTP endpoint:** `GET /metrics` — `/image/` ve `/healthz` dışında, çakışma yok.

**PR:**
```bash
git push -u origin feat/prometheus-metrics
gh pr create --base staging \
  --title "feat(server): add Prometheus metrics endpoint" \
  --body-file .github/PULL_REQUEST_TEMPLATE.md
```

---

#### Adım K — ADR-0006: Deployment Adapter Seçimi

> **Ön koşul:** Yok — Adım I ve J ile paralel açılabilir.

```bash
git checkout staging && git pull origin staging
git checkout -b docs/adr-0006-deployment-adapter
```

**ADR-0006 kapsamı:**
- Seçenekler: systemd-service packaging **veya** Fly.io adapter
- Karar kriterleri: ADR-0002 "deploys runtime vs. proxy" ayrımı, kompleksite
- V0.2 için **önerilen karar:** systemd packaging (sıfır dış platform bağımlılığı,
  `ROADMAP.md`'deki "standalone" önceliği)
- Fly.io: V0.2.x veya V0.3'e ertelenir

**Commit:**
```
docs(adr): add ADR-0006 for V0.2 deployment adapter selection
```

**PR:**
```bash
git push -u origin docs/adr-0006-deployment-adapter
gh pr create --base staging \
  --title "docs(adr): add ADR-0006 deployment adapter selection" \
  --body-file .github/PULL_REQUEST_TEMPLATE.md
```

---

#### Adım L — Systemd Deployment Adapter

> **Ön koşul:** ADR-0006 (Adım K) merge edilmiş olmalı.

```bash
git checkout staging && git pull origin staging
git checkout -b feat/adapter-systemd
```

**Kapsam:**
- `deploy/systemd/` — unit file template (`optivor.service`)
- GoReleaser'a `.deb` / `.rpm` package eklentisi (nfpm)
- `Makefile` hedefi: `make install` (systemd enable + start)
- `docs/deployment/systemd.md` — kurulum kılavuzu

**Commit sırası:**

```
feat(adapter/systemd): add systemd unit file template for optivor service

chore(goreleaser): add nfpm deb/rpm packaging for systemd deployment

feat(adapter/systemd): add Makefile install/uninstall targets

docs(deployment): add systemd deployment guide

test(adapter/systemd): add unit file validation smoke test
```

**PR:**
```bash
git push -u origin feat/adapter-systemd
gh pr create --base staging \
  --title "feat(adapter/systemd): add standalone systemd deployment adapter" \
  --body-file .github/PULL_REQUEST_TEMPLATE.md
```

---

#### Adım M — V0.2 E2E Kabul Testi Güncellemesi

> **Ön koşul:** I–L'nin tamamı merge edilmiş olmalı.

```bash
git checkout staging && git pull origin staging
git checkout -b test/v02-e2e-acceptance
```

**Commit:**
```
test(e2e): extend acceptance test for V0.2 features

Covers: /metrics endpoint, JSON log output, systemd unit file validity.
```

**PR:**
```bash
git push -u origin test/v02-e2e-acceptance
gh pr create --base staging \
  --title "test(e2e): V0.2 acceptance test suite" \
  --body-file .github/PULL_REQUEST_TEMPLATE.md
```

---

#### V0.2 Release

```bash
git checkout staging && git pull origin staging
gh pr create \
  --base main \
  --head staging \
  --title "release: promote staging → main (v0.2.0)" \
  --body "V0.2 milestone: structured logging, Prometheus metrics, systemd deployment adapter."

# Merge sonrası:
git checkout main && git pull origin main
git tag v0.2.0
git push origin v0.2.0
```

---

### V0.2 Kapsam Dışı

| Özellik | Hedef |
|---|---|
| CLI (`optivor init`, `optivor deploy`) | V0.3 |
| Fly.io / Cloudflare adapter | V0.3+ |
| OpenTelemetry tracing | V0.4 |
| `optivor doctor`, `optivor logs` CLI komutları | V0.4 |

---

## 📋 V0.3 — CLI

**ROADMAP.md hedefleri:**
- `optivor init`, `optivor deploy`, config scaffolding
- CLI, mevcut Deployment Adapter'ları (V0.2 itibarıyla systemd) orkestre eder

**ADR-0002 CLI katmanı sınırları:**
- CLI asla image processing yapmaz, storage'a doğrudan bağlanmaz
- CLI yalnızca runtime'ı ve deployment adapter'ları çağırır
- `cli/` paketi `internal/pipeline`, `internal/cache`'i import etmez

---

### V0.3 Adım Sırası

---

#### Adım N — ADR-0007: CLI Tasarım Kararları

> **Ön koşul:** V0.2 tüm PR'ları merge edilmiş, `main`'de `v0.2.0` tagı var.

```bash
git checkout staging && git pull origin staging
git checkout -b docs/adr-0007-cli-design
```

**ADR-0007 kapsamı:**
- CLI binary adı: `optivor` (tek binary, subcommand'larla)
- Framework seçimi: `cobra` (ADR-0001 "minimal dependencies" ile uyumlu)
- Config scaffolding: `optivor.yaml` şablonu
- Adapter discovery: `--adapter` flag veya config alanı
- `optivor init` çıktısı: dizin yapısı + `optivor.yaml` + `.gitignore` önerileri
- `optivor deploy` çıktısı: seçili adapter'ı çağırır (V0.2'de systemd)

**Commit:**
```
docs(adr): add ADR-0007 for CLI design and command structure
```

**PR:**
```bash
git push -u origin docs/adr-0007-cli-design
gh pr create --base staging \
  --title "docs(adr): add ADR-0007 CLI design decisions" \
  --body-file .github/PULL_REQUEST_TEMPLATE.md
```

---

#### Adım O — CLI İskeleti ve `optivor init`

> **Ön koşul:** ADR-0007 (Adım N) merge edilmiş olmalı.

```bash
git checkout staging && git pull origin staging
git checkout -b feat/cli-init
```

**Commit sırası:**

```
chore(deps): add cobra CLI framework dependency

feat(cli): add root command with version flag

feat(cli): implement 'optivor init' command with config scaffolding

test(cli): add table-driven tests for init command output

docs(readme): document 'optivor init' usage
```

**`optivor init` davranışı:**
- Çalışma dizinine `optivor.yaml` şablonu yazar (varsa üzerine yazma sorar)
- `.gitignore`'a `optivor.yaml`'daki secret alanlarını ekler
- Çıktı: kullanıcıya sonraki adımları gösteren kısa kılavuz

**PR:**
```bash
git push -u origin feat/cli-init
gh pr create --base staging \
  --title "feat(cli): add root command and 'optivor init' scaffolding" \
  --body-file .github/PULL_REQUEST_TEMPLATE.md
```

---

#### Adım P — `optivor deploy` Komutu

> **Ön koşul:** Adım O merge edilmiş olmalı.

```bash
git checkout staging && git pull origin staging
git checkout -b feat/cli-deploy
```

**Commit sırası:**

```
feat(cli): implement 'optivor deploy' command with adapter dispatch

feat(adapter/systemd): expose deploy entry-point callable from CLI

test(cli): add integration test for deploy command with systemd adapter

docs(readme): document 'optivor deploy' usage and adapter flag
```

**`optivor deploy` davranışı:**

| Flag | Default | Açıklama |
|---|---|---|
| `--adapter` | `systemd` (V0.3'te tek seçenek) | Hedef adapter |
| `--config` | `optivor.yaml` | Config dosyası yolu |
| `--dry-run` | false | Gerçek deploy yapmadan plan göster |

**Layer boundary:** `cli/` paketi `adapter/systemd`'yi out-of-process subprocess olarak çağırır — import yapmaz (ADR-0003 §Deployment Adapter).

**PR:**
```bash
git push -u origin feat/cli-deploy
gh pr create --base staging \
  --title "feat(cli): implement 'optivor deploy' with adapter dispatch" \
  --body-file .github/PULL_REQUEST_TEMPLATE.md
```

---

#### Adım Q — V0.3 E2E Kabul Testi

> **Ön koşul:** N–P'nin tamamı merge edilmiş olmalı.

```bash
git checkout staging && git pull origin staging
git checkout -b test/v03-e2e-acceptance
```

**Commit:**
```
test(e2e): extend acceptance test for V0.3 CLI features

Covers: 'optivor init' scaffolding output, 'optivor deploy --dry-run'
execution, and cobra help text completeness.
```

**Test senaryoları:**
- `optivor init` → `optivor.yaml` oluşturur, şema geçerliliği ✓
- `optivor deploy --dry-run --adapter systemd` → çıkış kodu 0, plan çıktısı ✓
- `optivor --version` → SemVer formatında çıktı ✓
- `optivor --help` → tüm subcommand'lar listelendi ✓

**PR:**
```bash
git push -u origin test/v03-e2e-acceptance
gh pr create --base staging \
  --title "test(e2e): V0.3 CLI acceptance test suite" \
  --body-file .github/PULL_REQUEST_TEMPLATE.md
```

---

#### V0.3 Release

```bash
git checkout staging && git pull origin staging
gh pr create \
  --base main \
  --head staging \
  --title "release: promote staging → main (v0.3.0)" \
  --body "V0.3 milestone: CLI (optivor init, optivor deploy), systemd adapter integration."

# Merge sonrası:
git checkout main && git pull origin main
git tag v0.3.0
git push origin v0.3.0
```

---

### V0.3 Kapsam Dışı

| Özellik | Hedef |
|---|---|
| `optivor doctor`, `optivor logs`, `optivor metrics` | V0.4 |
| OpenTelemetry tracing | V0.4 |
| Scope/path bazlı auth policy | V0.3+ (ADR gerekli) |
| Fly.io / Cloudflare adapter | V1 |
| Distributed rate limiting (Redis) | Belirsiz |

---

## 📋 V0.4 — Observability

**ROADMAP.md hedefleri:**
- OpenTelemetry tracing
- `optivor doctor`, `optivor logs`, `optivor metrics` diagnostics komutları

**Not:** V0.2'de Prometheus metrics zaten eklendi (`/metrics` endpoint). V0.4,
bunu OpenTelemetry distributed tracing ile tamamlar ve CLI diagnostics
komutlarıyla kullanıcıya sunar.

---

### V0.4 Adım Sırası

---

#### Adım R — ADR-0008: OpenTelemetry Entegrasyonu

> **Ön koşul:** V0.3 tüm PR'ları merge edilmiş, `main`'de `v0.3.0` tagı var.

```bash
git checkout staging && git pull origin staging
git checkout -b docs/adr-0008-opentelemetry
```

**ADR-0008 kapsamı:**
- SDK seçimi: `go.opentelemetry.io/otel` (CNCF standard)
- Exporter seçimi: OTLP/gRPC (default) + stdout (geliştirme)
- Propagation: W3C TraceContext + Baggage
- Instrumentation noktaları: HTTP handler, pipeline transform, storage driver, cache
- Sampling stratejisi: V0.4'te `AlwaysSample`; üretim için `TraceIDRatioBased` önerilir
- ADR-0002 layer sınırı: trace context `internal/server`'dan `internal/pipeline`'a
  context.Context üzerinden geçer — OTel import'ları pipeline'da kabul edilir

**Commit:**
```
docs(adr): add ADR-0008 for OpenTelemetry tracing integration
```

**PR:**
```bash
git push -u origin docs/adr-0008-opentelemetry
gh pr create --base staging \
  --title "docs(adr): add ADR-0008 OpenTelemetry tracing decisions" \
  --body-file .github/PULL_REQUEST_TEMPLATE.md
```

---

#### Adım S — OpenTelemetry Tracing İmplementasyonu

> **Ön koşul:** ADR-0008 (Adım R) merge edilmiş olmalı.

```bash
git checkout staging && git pull origin staging
git checkout -b feat/otel-tracing
```

**Commit sırası:**

```
chore(deps): add opentelemetry-go SDK and OTLP exporter dependencies

feat(server): initialize OTel TracerProvider with configurable OTLP endpoint

feat(server): instrument HTTP handler with OTel span (method, route, status)

feat(pipeline): propagate trace context through transform operations

feat(driver/s3): add OTel span for storage GetObject calls

feat(cache/fs): add OTel span for cache hit/miss operations

test(server): add trace propagation integration test

docs(readme): document OTLP endpoint configuration and local Jaeger setup
```

**Config değişikliği:**
```yaml
telemetry:
  enabled: true
  otlp_endpoint: ""        # boşsa stdout exporter kullanılır
  service_name: "optivor"
  sampling_ratio: 1.0      # 1.0 = AlwaysSample
```

**Span hiyerarşisi:**
```
HTTP request span (server)
  └── pipeline.Transform span
        ├── storage.GetObject span
        └── cache.Lookup span
```

**PR:**
```bash
git push -u origin feat/otel-tracing
gh pr create --base staging \
  --title "feat(server,pipeline): add OpenTelemetry distributed tracing" \
  --body-file .github/PULL_REQUEST_TEMPLATE.md
```

---

#### Adım T — `optivor doctor` Komutu

> **Ön koşul:** Adım S merge edilmiş olmalı.

```bash
git checkout staging && git pull origin staging
git checkout -b feat/cli-doctor
```

**`optivor doctor` kontrolleri:**

| Kontrol | Başarı | Hata |
|---|---|---|
| Config dosyası parse edilebiliyor | ✅ | ❌ + satır numarası |
| S3 bucket erişilebilir (ListObjects) | ✅ | ❌ + hata mesajı |
| `OPTIVOR_AUTH_SECRET` set (enabled ise) | ✅ | ⚠️ warning |
| Runtime binary çalıştırılabilir | ✅ | ❌ |
| libvips sürümü ≥ minimum | ✅ | ❌ + kurulum önerisi |

**Commit sırası:**

```
feat(cli): implement 'optivor doctor' health check command

feat(cli/doctor): add config parse check

feat(cli/doctor): add S3 connectivity check

feat(cli/doctor): add auth secret presence check

feat(cli/doctor): add libvips version check

test(cli): add doctor command unit tests with mock checks

docs(readme): document 'optivor doctor' usage
```

**PR:**
```bash
git push -u origin feat/cli-doctor
gh pr create --base staging \
  --title "feat(cli): add 'optivor doctor' diagnostics command" \
  --body-file .github/PULL_REQUEST_TEMPLATE.md
```

---

#### Adım U — `optivor logs` ve `optivor metrics` Komutları

> **Ön koşul:** Adım T merge edilmiş olmalı.

```bash
git checkout staging && git pull origin staging
git checkout -b feat/cli-logs-metrics
```

**`optivor logs` davranışı:**
- Çalışan optivor process'inin log dosyasını tail eder (systemd: `journalctl -u optivor -f`)
- `--since`, `--until`, `--lines` flag'leri
- JSON log formatını tablo olarak pretty-print edebilir (`--format table`)

**`optivor metrics` davranışı:**
- Çalışan runtime'ın `/metrics` endpoint'ine bağlanır
- Seçili metrikleri (`optivor_requests_total`, latency, cache hit ratio) özet tablo olarak gösterir
- `--watch` flag: her N saniyede bir yeniler

**Commit sırası:**

```
feat(cli): implement 'optivor logs' command with journalctl integration

feat(cli): implement 'optivor metrics' command with /metrics scraping

test(cli): add logs and metrics command tests with mock adapter

docs(readme): document 'optivor logs' and 'optivor metrics' usage
```

**PR:**
```bash
git push -u origin feat/cli-logs-metrics
gh pr create --base staging \
  --title "feat(cli): add 'optivor logs' and 'optivor metrics' diagnostics commands" \
  --body-file .github/PULL_REQUEST_TEMPLATE.md
```

---

#### Adım V — V0.4 E2E Kabul Testi

> **Ön koşul:** R–U'nun tamamı merge edilmiş olmalı.

```bash
git checkout staging && git pull origin staging
git checkout -b test/v04-e2e-acceptance
```

**Commit:**
```
test(e2e): extend acceptance test for V0.4 observability features

Covers: OTel span export to stdout collector, doctor command exit codes,
metrics command output format, logs command --lines flag.
```

**Test senaryoları:**
- Image request → stdout exporter'da span görünür (service.name=optivor) ✓
- `optivor doctor` geçerli config ile → çıkış kodu 0 ✓
- `optivor doctor` geçersiz S3 config ile → çıkış kodu 1, hata mesajı ✓
- `optivor metrics` → `optivor_requests_total` satırı çıktıda ✓
- `optivor logs --lines 10` → 10 satır döner, hata yok ✓

**PR:**
```bash
git push -u origin test/v04-e2e-acceptance
gh pr create --base staging \
  --title "test(e2e): V0.4 observability acceptance test suite" \
  --body-file .github/PULL_REQUEST_TEMPLATE.md
```

---

#### V0.4 Release

```bash
git checkout staging && git pull origin staging
gh pr create \
  --base main \
  --head staging \
  --title "release: promote staging → main (v0.4.0)" \
  --body "V0.4 milestone: OpenTelemetry tracing, optivor doctor/logs/metrics diagnostics CLI."

# Merge sonrası:
git checkout main && git pull origin main
git tag v0.4.0
git push origin v0.4.0
```

---

### V0.4 Kapsam Dışı

| Özellik | Hedef |
|---|---|
| Storage Driver interface finalize / dış katkı dokümantasyonu | V1 |
| Runtime Module mechanism ADR | V1 |
| Ek Deployment Adapter'lar (Fly.io, Cloudflare, AWS, Kubernetes) | V1 |
| Dashboard / web UI | Planlanmadı |
| AI-based transformations | Planlanmadı |
| Multi-node / horizontally scaled runtime | Planlanmadı |

---

## Genel Branch & Commit Kuralları (AGENTS.md özeti)

> Detay için her zaman `docs/agents/AGENTS.md`'e bak. Bu özet hatırlatma içindir.

### Branch modeli

```
feature branches  →  staging  →  main
   (PRs land here)   (integration)  (releases only)
```

- `main` = sadece release tag'larında güncellenir
- `staging` = tüm feature PR'larının hedefi
- Hiçbir şey doğrudan `main`'e commit edilmez

### Commit formatı (Conventional Commits)

```
<type>(<scope>): <özet, emir kipi, nokta yok>

<isteğe bağlı gövde: ne değişti ve neden>

<isteğe bağlı: Refs #123, See ADR-0005>
```

Type'lar: `feat`, `fix`, `docs`, `refactor`, `test`, `chore`, `revert`

Scope örnekleri: `server`, `pipeline`, `cache/fs`, `driver/s3`, `adapter/systemd`, `cli`, `adr`

### PR açma

```bash
git push -u origin <branch>
gh pr create \
  --base staging \
  --title "<type>(<scope>): <özet>" \
  --body-file .github/PULL_REQUEST_TEMPLATE.md
```

- Draft PR büyük işler için: `gh pr create --draft ...`
- CI yeşil olmadan review istenmez
- Kendi PR'ını merge etme — CODEOWNERS bekle

### Rebase kuralı

```bash
git checkout staging && git pull origin staging
git checkout <branch>
git rebase staging   # merge değil, rebase
```

---

## Mimari Kararlar (ADR referansları)

| ADR | Konu | Durum |
|---|---|---|
| ADR-0000 | Project Scope | ✅ Accepted |
| ADR-0001 | Project Philosophy | ✅ Accepted |
| ADR-0002 | High-Level Architecture | ✅ Accepted |
| ADR-0003 | Extensibility Model | ✅ Accepted |
| ADR-0004 | Technology Choices | ✅ Accepted |
| ADR-0005 | Signed URL & Auth | ✅ V0.1 (tamamlandı) |
| ADR-0006 | Deployment Adapter | 📋 V0.2 (Adım K) |
| ADR-0007 | CLI Design | 📋 V0.3 (Adım N) |
| ADR-0008 | OpenTelemetry Tracing | 📋 V0.4 (Adım R) |

---

## Referanslar

- [ROADMAP.md](../../ROADMAP.md)
- [AGENTS.md](AGENTS.md)
- [CICD.md](CICD.md)
- [ADR-0000](../adr/0000-project-scope.md)
- [ADR-0001](../adr/0001-project-philosophy.md)
- [ADR-0002](../adr/0002-high-level-architecture.md)
- [ADR-0003](../adr/0003-extensibility-model.md)
- [ADR-0004](../adr/0004-technology-choices.md)
- [Architecture: runtime.md](../architecture/runtime.md)
