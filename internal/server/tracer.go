package server

import (
	"context"
	"fmt"
	"io"

	"github.com/optivor/optivor/internal/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

// InitTracer initializes OpenTelemetry TracerProvider based on configuration.
func InitTracer(cfg *config.Config, out io.Writer) (*sdktrace.TracerProvider, error) {
	if !cfg.Telemetry.Enabled {
		tp := sdktrace.NewTracerProvider()
		otel.SetTracerProvider(tp)
		return tp, nil
	}

	exporter, err := stdouttrace.New(stdouttrace.WithWriter(out))
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout trace exporter: %w", err)
	}

	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.Telemetry.ServiceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create otel resource: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(cfg.Telemetry.SamplingRatio))),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tp, nil
}

// Tracer returns named tracer instance.
func Tracer() trace.Tracer {
	return otel.GetTracerProvider().Tracer("optivor")
}
