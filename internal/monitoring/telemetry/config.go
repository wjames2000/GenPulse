package telemetry

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

type Config struct {
	Enabled            bool
	ServiceName        string
	ServiceVersion     string
	Environment        string
	OTLPEndpoint       string
	OTLPInsecure       bool
	ConsoleExport      bool
	BatchTimeout       time.Duration
	ExportTimeout      time.Duration
	MaxExportBatchSize int
}

func DefaultConfig() Config {
	return Config{
		Enabled:            true,
		ServiceName:        "genpulse",
		ServiceVersion:     "1.0.0",
		Environment:        "development",
		OTLPEndpoint:       "localhost:4317",
		OTLPInsecure:       true,
		ConsoleExport:      true,
		BatchTimeout:       5 * time.Second,
		ExportTimeout:      30 * time.Second,
		MaxExportBatchSize: 512,
	}
}

type Telemetry struct {
	tracerProvider *sdktrace.TracerProvider
	config         Config
}

func NewTelemetry(config Config) (*Telemetry, error) {
	t := &Telemetry{
		config: config,
	}

	if !config.Enabled {
		return t, nil
	}

	var exporters []sdktrace.SpanExporter

	if config.ConsoleExport {
		consoleExporter, err := stdouttrace.New(
			stdouttrace.WithWriter(os.Stdout),
			stdouttrace.WithPrettyPrint(),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create console exporter: %w", err)
		}
		exporters = append(exporters, consoleExporter)
	}

	if config.OTLPEndpoint != "" {
		ctx, cancel := context.WithTimeout(context.Background(), config.ExportTimeout)
		defer cancel()

		otlpExporter, err := otlptracegrpc.New(ctx,
			otlptracegrpc.WithEndpoint(config.OTLPEndpoint),
			otlptracegrpc.WithInsecure(),
		)
		if err != nil {
			log.Printf("Warning: failed to create OTLP exporter: %v", err)
		} else {
			exporters = append(exporters, otlpExporter)
		}
	}

	if len(exporters) == 0 {
		return t, nil
	}

	var spanProcessor sdktrace.SpanProcessor
	if len(exporters) == 1 {
		spanProcessor = sdktrace.NewBatchSpanProcessor(exporters[0],
			sdktrace.WithBatchTimeout(config.BatchTimeout),
			sdktrace.WithMaxExportBatchSize(config.MaxExportBatchSize),
			sdktrace.WithMaxQueueSize(1000),
		)
	} else {
		spanProcessor = sdktrace.NewBatchSpanProcessor(exporters[0],
			sdktrace.WithBatchTimeout(config.BatchTimeout),
			sdktrace.WithMaxExportBatchSize(config.MaxExportBatchSize),
			sdktrace.WithMaxQueueSize(1000),
		)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(spanProcessor),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(config.ServiceName),
			semconv.ServiceVersion(config.ServiceVersion),
			attribute.String("environment", config.Environment),
			attribute.String("application", "genpulse"),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	t.tracerProvider = tp
	return t, nil
}

func (t *Telemetry) Shutdown(ctx context.Context) error {
	if t.tracerProvider != nil {
		return t.tracerProvider.Shutdown(ctx)
	}
	return nil
}

func (t *Telemetry) Tracer(name string) interface{} {
	if t.tracerProvider == nil {
		return otel.GetTracerProvider().Tracer(name)
	}
	return t.tracerProvider.Tracer(name)
}

func (t *Telemetry) IsEnabled() bool {
	return t.config.Enabled && t.tracerProvider != nil
}
