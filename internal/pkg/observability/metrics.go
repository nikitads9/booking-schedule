package observability

import (
	"booking-schedule/internal/logger/sl"
	"context"
	"log/slog"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/cpu"

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
		return nil, err
	}

	exporter, err := prometheus.New()
	if err != nil {
		return nil, err
	}

	provider := metricsdk.NewMeterProvider(
		metricsdk.WithResource(resource),
		metricsdk.WithReader(exporter),
	)

	otel.SetMeterProvider(provider)

	return provider.Meter(svcName), nil
}

func CollectMachineResourceMetrics(meter metric.Meter, logger *slog.Logger) {
	const op = "observability.CollectMachineResourceMetrics"

	log := logger.With(
		slog.String("op", op))

	period := 10 * time.Second
	ticker := time.NewTicker(period)

	var Mb uint64 = 1_048_576 // number of bytes in a MB
	for {
		<-ticker.C
		// This will be executed every "period" of time passes
		_, err := meter.Float64ObservableGauge(
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
		if err != nil {
			log.Error("could not collect allocated memory gauge", sl.Err(err))
		}

		_, err = meter.Float64ObservableGauge(
			"process.cpu_usage",
			metric.WithFloat64Callback(
				func(ctx context.Context, fo metric.Float64Observer) error {
					cpuUsage, err := cpu.Percent(0, false)
					if err != nil {
						return err
					}

					if cpuUsage != nil {
						fo.Observe(cpuUsage[0])
					}

					return nil
				},
			),
			metric.WithUnit("%"),
		)
		if err != nil {
			log.Error("could not collect cpu usage", sl.Err(err))
		}
	}

}
