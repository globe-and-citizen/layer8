package opentelemetry

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	metricsdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// NewMeter creates a new metric.Meter that can create any metric reporter
// you might want to use in your application.
func NewMeter(ctx context.Context) error {
	provider, err := newMeterProvider(ctx)
	if err != nil {
		return fmt.Errorf("could not create meter provider: %w", err)
	}

	otel.SetMeterProvider(provider)

	return nil
}

// newMeterProvcider initialize the application resource, connects to the
// OpenTelemetry Collector and configures the metric poller that will be used
// to collect the metrics and send them to the OpenTelemetry Collector.
func newMeterProvider(ctx context.Context) (metric.MeterProvider, error) {
	// Interval which the metrics will be reported to the collector
	interval := 10 * time.Second

	resource, err := getResource()
	if err != nil {
		return nil, fmt.Errorf("could not get resource: %w", err)
	}

	collectorExporter, err := getOtelMetricsCollectorExporter(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get collector exporter: %w", err)
	}

	periodicReader := metricsdk.NewPeriodicReader(collectorExporter,
		metricsdk.WithInterval(interval),
	)

	provider := metricsdk.NewMeterProvider(
		metricsdk.WithResource(resource),
		metricsdk.WithReader(periodicReader),
	)

	return provider, nil
}

// getResource creates the resource that describes our application.
func getResource() (*resource.Resource, error) {
	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("layer8"),
	)

	return resource, nil
}

// getOtelMetricsCollectorExporter creates a metric exporter
func getOtelMetricsCollectorExporter(ctx context.Context) (metricsdk.Exporter, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	exporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")),
		otlpmetricgrpc.WithCompressor("gzip"),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("could not create metric exporter: %w", err)
	}

	return exporter, nil
}
