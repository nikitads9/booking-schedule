package observability

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"go.opentelemetry.io/otel/exporters/prometheus"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	metricsdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func NewMeter(ctx context.Context, svcName string) (metric.Meter, error) {
	resource, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(svcName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("could not merge resources: %w", err)
	}

	exporter, err := prometheus.New()
	if err != nil {
		return nil, fmt.Errorf("could not get collector exporter: %w", err)
	}

	provider := metricsdk.NewMeterProvider(
		metricsdk.WithResource(resource),
		metricsdk.WithReader(exporter),
	)

	otel.SetMeterProvider(provider)

	return provider.Meter(svcName), nil
}

func CollectMachineResourceMetrics(meter metric.Meter) {
	period := 5 * time.Second
	ticker := time.NewTicker(period)

	var Mb uint64 = 1_048_576 // number of bytes in a MB

	for {
		<-ticker.C
		// This will be executed every "period" of time passes
		meter.Float64ObservableGauge(
			"process.allocated_memory",
			metric.WithFloat64Callback(
				func(ctx context.Context, fo metric.Float64Observer) error {
					var memStats runtime.MemStats
					runtime.ReadMemStats(&memStats)

					allocatedMemoryInMB := float64(memStats.Alloc) / float64(Mb)
					fo.Observe(allocatedMemoryInMB)
					return nil
				},
			),
		)
	}
}
