package trace

import (
	"context"
	"fmt"

	jaegerPropagator "go.opentelemetry.io/contrib/propagators/jaeger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func NewTracer(ctx context.Context, jaegerEndpoint string, svcName string, samplingRate float64) (trace.Tracer, error) {
	// create jaeger exporter
	exporter, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpointURL(jaegerEndpoint), otlptracehttp.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("unable to initialize exporter due: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		//sdktrace.WithSampler(trace.),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(samplingRate)),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(svcName),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(jaegerPropagator.Jaeger{})

	// returns tracer
	return otel.Tracer(svcName), nil
}
