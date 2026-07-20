# ADR-0008: OpenTelemetry Tracing Integration

## Status

Accepted

## Context

As Optivor evolves into V0.4 (Observability), developers need distributed tracing to observe end-to-end latency across HTTP handlers, image transformation pipelines, storage retrieval operations, and filesystem cache lookups.

Per **ADR-0002** (High-Level Architecture), core runtime packages like `internal/pipeline` and `internal/storage` must remain decoupled from specific cloud vendors or monolithic observability platforms.

## Decision

1. **OpenTelemetry SDK Standard**:
   Adopt OpenTelemetry Go (`go.opentelemetry.io/otel`) as the CNCF-standard distributed tracing API and SDK.

2. **Exporters**:
   Support OTLP/gRPC exporter for production collectors (e.g. Jaeger, Grafana Tempo, OpenTelemetry Collector) and a stdout exporter for local development.

3. **Context Propagation**:
   Enforce W3C TraceContext propagation across HTTP request boundaries. Trace context flows through standard `context.Context` parameters down to pipeline, storage, and cache calls.

4. **Span Hierarchy**:
   - `HTTP GET /image/*`: Root span (`server`)
   - `pipeline.Transform`: Child span (`pipeline`)
   - `storage.GetObject`: Child span (`storage`)
   - `cache.Lookup`: Child span (`cache`)

5. **Configuration Schema**:
   Managed under `telemetry:` in `optivor.yaml`:
   ```yaml
   telemetry:
     enabled: true
     otlp_endpoint: "" # default to stdout exporter if empty
     service_name: "optivor"
     sampling_ratio: 1.0
   ```

## Consequences

- End-to-end request latency and bottlenecks across storage, caching, and libvips image processing become transparently observable.
- Zero vendor lock-in for observability platforms since OpenTelemetry exports via standard OTLP protocol.
