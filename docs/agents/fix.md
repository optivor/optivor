# Optivor v0.4.0 — Sistem Denetimi ve Düzeltme Planı

> **Tarih:** 2026-07-20  
> **Mevcut versiyon:** v0.4.0  
> **Denetim kapsamı:** ADR bütünlüğü, kod kalitesi, CI/CD, README, contributor deneyimi, güvenlik, versiyon planlaması

---

## Özet Bulgular

| Kategori | Ciddiyet | Adet |
|---|---|---|
| Kritik kod hatası | 🔴 Yüksek | 3 |
| CI/CD eksikliği | 🟠 Orta | 4 |
| ADR tutarsızlığı | 🟠 Orta | 3 |
| Contributor deneyimi | 🟡 Düşük | 4 |
| README / dokümantasyon | 🟡 Düşük | 5 |

---

## Versiyon Planlaması

Yeni **feature** yok — tüm düzeltmeler patch sürümleri:

```
v0.4.0  (mevcut)
  │
  ├── v0.4.1  — Kritik kod düzeltmeleri (3 fix, güvenlik dahil)
  ├── v0.4.2  — CI/CD ve CODEOWNERS düzeltmeleri
  └── v0.4.3  — Dokümantasyon, README, contributor UX iyileştirmeleri
```

> **V1 yolu:** v0.4.3 stable → `staging → main` release PR → `v0.5.0`'a
> geçilir (Storage Driver interface finalize, Runtime Module ADR — bkz. ROADMAP.md).

---

## 🔴 Kritik — v0.4.1

### FIX-01: `cache/fs` — OTel span bağlamı kopuk (orphan spans)

**Dosya:** `internal/cache/fs/fscache.go` satır 40, 74

**Sorun:** Her iki OTel span `context.Background()` ile açılıyor. Bu gelen HTTP
request'in trace context'ini koparır. ADR-0008'in span hiyerarşisi
(`server → pipeline → cache.Lookup`) hiçbir zaman oluşmaz; tüm cache span'ları
orphan span olarak kayıt altına alınır.

```go
// YANLIŞ (şu anki durum)
_, span := otel.Tracer("optivor").Start(context.Background(), "cache.Get")

// DOĞRU — ctx parametresi eklenmeli
func (c *FSCache) Get(ctx context.Context, key string, params pipeline.TransformParams) ([]byte, string, bool, error) {
    _, span := otel.Tracer("optivor").Start(ctx, "cache.Get")
```

`internal/cache/cache.go` `Cache` interface'i de `ctx context.Context`
parametresini alacak şekilde güncellenmeli. Bu `server.go` içindeki çağrılara yansır.

**ADR referansı:** ADR-0008 §3 — "trace context flows through `context.Context`
parameters down to pipeline, storage, and cache calls."

---

### FIX-02: `cli/doctor.go` — libvips `Startup` çağrısı `Shutdown` olmadan

**Dosya:** `internal/cli/doctor.go` satır 68

**Sorun:** `vips.Startup(nil)` yapılıyor, `defer vips.Shutdown()` yok. libvips
bir C kütüphanesi; Startup ile tahsis edilen kaynaklar (thread pool, bellek)
process ömrü boyunca açık kalır. Kısa ömürlü CLI komutu için kaynak sızıntısı.

Ayrıca `doctor.go` sadece libvips'i başlatıp başarılı sayıyor; gerçek **sürüm
kontrolü** yok. `plan.md` Adım T "libvips sürümü ≥ minimum" kontrolünü
tanımlıyor ancak implementasyon bunu yapmıyor.

```go
// Düzeltme
vips.Startup(nil)
defer vips.Shutdown()
info := vips.GetVersionString()
fmt.Printf("  ✅ libvips %s initialized\n", info)
```

---

### FIX-03: `server/auth.go` — ADR-0005 ile imza girdisi tutarsızlığı

**Dosya:** `internal/server/auth.go`

**Sorun:** ADR-0005 §1 imza girdisini şöyle tanımlıyor:
```
path + "?expires=" + expires
```

Ancak mevcut `GenerateSignature` fonksiyonu:
```go
sigInput := path
if len(q) > 0 {
    sigInput += "?" + q.Encode()  // expires + w + h + fit + format hepsini ekliyor
}
```

Bu tutarsızlık istemci SDK'larının doğru imza üretmesini imkânsız kılar. Harici
katkıcılar ADR-0005'i referans alarak yazacakları istemci kodları çalışmayacak.

**Öneri:** Mevcut implementasyonu koru (tüm param'ları imzalamak daha güvenlidir),
ADR-0005'i güncelleyerek gerçek davranışı belgele. `GenerateSignature` fonksiyonuna
godoc yorum ekle.

---

## 🟠 Orta — v0.4.2

### FIX-04: CI workflow — `golangci-lint` ve `govulncheck` eksik

**Dosya:** `.github/workflows/ci.yml`

CI workflow sadece `go build`, `go vet`, `go test` içeriyor. `docs/agents/CICD.md`
§2'de zorunlu olarak tanımlanan `golangci-lint` ve `govulncheck` hiç çalışmıyor.

```yaml
# Eklenecek adımlar
- name: Run golangci-lint
  uses: golangci/golangci-lint-action@v4
  with:
    version: latest

- name: Install govulncheck
  run: go install golang.org/x/vuln/cmd/govulncheck@latest

- name: Run govulncheck
  run: govulncheck ./...
```

CICD.md'de ayrıca tanımlanan ama henüz implement edilmemiş kontroller:
- Layer boundary check (ADR-0002 enforcement)
- Commit message lint (Conventional Commits)
- GoReleaser snapshot dry-run (release PR'larda)

---

### FIX-05: CI vs `go.mod` versiyon uyumsuzluğu

**Dosya:** `.github/workflows/ci.yml` satır 20 + `go.mod` satır 3

CI `go-version: '1.23'` kullanıyor; `go.mod` `go 1.25.0` gerektiriyor.

```yaml
- uses: actions/setup-go@v5
  with:
    go-version: '1.25'
```

---

### FIX-06: `CODEOWNERS` — Gerçek dosya yolları ile eşleşmiyor

**Dosya:** `CODEOWNERS`

CODEOWNERS `/runtime/` ve `/cli/` yollarını tanımlıyor; bu dizinler yok.
GitHub var olmayan yolları **sessizce yok sayar** — hiçbir PR'a otomatik reviewer atanmıyor.

```
# Doğru yollar
/internal/server/     @OPTIVOR_ORG/core-maintainers
/internal/pipeline/   @OPTIVOR_ORG/core-maintainers
/internal/storage/    @OPTIVOR_ORG/core-maintainers
/internal/cache/      @OPTIVOR_ORG/core-maintainers
/internal/cli/        @OPTIVOR_ORG/core-maintainers
/internal/config/     @OPTIVOR_ORG/core-maintainers
/cmd/                 @OPTIVOR_ORG/core-maintainers
/docs/adr/            @OPTIVOR_ORG/core-maintainers
```

---

### FIX-07: `tracer.go` — OTLP/gRPC exporter implement edilmemiş

**Dosya:** `internal/server/tracer.go`

ADR-0008 §2 iki exporter tanımlıyor: OTLP/gRPC (production) ve stdout (dev).
Mevcut kod `cfg.Telemetry.OTLPEndpoint` alanını tamamen yok sayıyor; dolu bir
değer girilse de her zaman stdout exporter kullanılıyor.

```go
// Düzeltme — OTLPEndpoint doluysa gRPC exporter kullan
if cfg.Telemetry.OTLPEndpoint != "" {
    // go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc
    conn, _ := grpc.NewClient(cfg.Telemetry.OTLPEndpoint,
        grpc.WithTransportCredentials(insecure.NewCredentials()))
    exporter, err = otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
}
```

`go.mod`'a `go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc`
bağımlılığı eklenmeli.

---

## 🟡 Düşük — v0.4.3

### FIX-08: README — Versiyon badge yok, AVIF bağımlılığı eksik, OTel bölümü yok

**Dosya:** `README.md`

1. Version badge yok; hangi sürümü indireceği belli değil.
2. `apt install libvips-dev` yazıyor; AVIF için `libheif-dev` de gerekli.
3. OpenTelemetry yapılandırması için herhangi bir bölüm yok (v0.4.0 ile eklendi).
4. "Architecture Overview" diyagramı cache ve pipeline katmanlarını göstermiyor.

```markdown
[![Release](https://img.shields.io/github/v/release/optivor/optivor)](...)
[![Go Report Card](https://goreportcard.com/badge/github.com/optivor/optivor)](...)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue)](LICENSE)

# Quick Start — build
apt install libvips-dev libheif-dev   # AVIF için libheif-dev gerekli
```

---

### FIX-09: ROADMAP.md — Tamamlanan milestone'lar "in progress" görünüyor

**Dosya:** `ROADMAP.md`

V0–V0.4 tamamlanmış; ROADMAP statüleri güncellenmeli:

```markdown
## V0 — Prove the core loop ✅ Released as v0.1.0-alpha
## V0.1 — Make it safe to expose ✅ Released as v0.1.0
## V0.2 — First deployment adapter ✅ Released as v0.2.0
## V0.3 — CLI ✅ Released as v0.3.0
## V0.4 — Observability ✅ Released as v0.4.0
```

---

### FIX-10: `docs/agents/plan.md` — ADR durumları eski

**Dosya:** `docs/agents/plan.md` satır 1171–1173

ADR-0006, ADR-0007, ADR-0008 hepsi `Accepted` durumunda ama plan.md bunları
hâlâ 📋 (yapılacak) olarak gösteriyor. Tüm milestone sütunları ✅ olarak
güncellenmeli.

---

### FIX-11: `SECURITY.md` eksik

**Eksik:** `SECURITY.md`

GitHub, bu dosya olmadan "Report a vulnerability" butonunu göstermiyor.
ADR-0001'in "contributor-first" taahhütüyle çelişiyor.

```markdown
# Security Policy

## Supported Versions
| Version | Supported |
|---------|-----------|
| 0.4.x   | ✅ |
| < 0.4   | ❌ |

## Reporting a Vulnerability
GitHub private vulnerability reporting kullanın veya
security@optivor.io adresine e-posta gönderin.
72 saat içinde yanıt, 14 gün içinde patch hedefliyoruz.
```

---

### FIX-12: `CONTRIBUTING.md` — Yerel geliştirme kurulumu eksik

**Dosya:** `CONTRIBUTING.md`

MinIO ile yerel test, `make test-e2e` kullanımı, gerekli sistem bağımlılıkları
(`libvips-dev`, `libheif-dev`) dokümante edilmemiş. İlk katkıcı için bu bilgiler
kritik.

---

## ADR Doğruluk Özeti

| ADR | Durum | Sorun |
|---|---|---|
| ADR-0000: Project Scope | ✅ | — |
| ADR-0001: Philosophy | ✅ | — |
| ADR-0002: Architecture | ✅ | Layer boundary check CI'da eksik (FIX-04) |
| ADR-0003: Extensibility | ✅ | Runtime Module ADR erteleme belgelenmiş |
| ADR-0004: Tech Choices | ⚠️ | Go version satırı 1.23 → 1.25 güncellenmeli |
| ADR-0005: Signed URL | ⚠️ | İmza girdisi tanımı impl. ile eşleşmiyor (FIX-03) |
| ADR-0006: Deployment | ✅ | — |
| ADR-0007: CLI | ✅ | — |
| ADR-0008: OTel | ⚠️ | OTLP/gRPC eksik (FIX-07), context prop. kırık (FIX-01) |

---

## Contributor Deneyimi Değerlendirmesi

### İyi Olan ✅
- CONTRIBUTING.md mimari katmanları net açıklıyor
- Issue template'leri (bug, feature, driver/adapter) scope disiplinini zorluyor
- PR template ADR referansı istiyor
- AGENTS.md branch/commit/PR kurallarını AI agent ve insanlar için net tanımlamış
- ADR-0000–0003 "neden?" sorularını önceden cevaplayan kurumsal kalitede belgeler

### Düzeltilmesi Gereken 🔧
1. CODEOWNERS yolları yanlış → hiçbir PR'a otomatik reviewer atanmıyor (FIX-06)
2. Yerel geliştirme talimatları yok → ilk katkı engeli (FIX-12)
3. SECURITY.md yok → güvenlik açığı bildirme yolu belirsiz (FIX-11)
4. CI eksiklikleri → golangci-lint, govulncheck çalışmıyor (FIX-04)

---

## Uygulama Planı (Branch Sırası)

### v0.4.1
```bash
fix/cache-otel-context-propagation   # FIX-01
fix/doctor-libvips-shutdown          # FIX-02
fix/auth-signature-adr-alignment     # FIX-03 + ADR-0005 güncelleme
```

### v0.4.2
```bash
chore/ci-workflow-lint-govulncheck   # FIX-04 + FIX-05
chore/codeowners-fix-paths           # FIX-06
feat/otel-otlp-grpc-exporter         # FIX-07
```

### v0.4.3
```bash
docs/readme-v040-update              # FIX-08
docs/roadmap-status-update           # FIX-09 + FIX-10
docs/add-security-policy             # FIX-11
docs/contributing-local-dev-setup    # FIX-12
```

---

## Sonraki Major Milestone: V1 Hazırlığı

| Konu | Gereksinim |
|---|---|
| Storage Driver interface finalize | ADR-0009 gerekli |
| Runtime Module mechanism (WASM vs in-process) | ADR-0010 gerekli |
| Ek Deployment Adapter'lar (Fly.io, Kubernetes) | ADR-0006 genişletme |
| Homebrew tap / APT / RPM dağıtım kanalları | GoReleaser config |

---

*Bu belge `v0.4.0` codebase'i üzerinde yapılan tam sistem denetiminin çıktısıdır.
Her fix kendi feature branch'i ve PR'ı ile uygulanmalı; AGENTS.md §1–§5 sırasına uyulmalıdır.*
