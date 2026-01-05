// otel is a package for setting up and initializing OpenTelemetry SDK components. Plumbing only.
package otel

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

// setupOTelSDK bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func SetupOTelSDK(ctx context.Context, collectorEndpoint, serviceName string) (func(context.Context) error, error) {
	var shutdownFuncs []func(context.Context) error
	var err error

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown := func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}

		shutdownFuncs = nil
		return err
	}

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) { err = errors.Join(inErr, shutdown(ctx)) }

	// set up propagator.
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	// set up global trace provider.
	tracerProvider, err := newOTLPTracerProvider(collectorEndpoint, serviceName)
	if err != nil {
		handleErr(err)
		return shutdown, err
	}

	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	// set up meter provider.
	meterProvider, err := newOTLPMeterProvider(collectorEndpoint, serviceName)
	if err != nil {
		handleErr(err)
		return shutdown, err
	}

	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	// set up logger provider.
	loggerProvider, err := newOTLPLoggerProvider(collectorEndpoint, serviceName)
	if err != nil {
		handleErr(err)
		return shutdown, err
	}
	shutdownFuncs = append(shutdownFuncs, loggerProvider.Shutdown)

	global.SetLoggerProvider(loggerProvider)

	return shutdown, err
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newOTLPTracerProvider(collectorEndpoint, serviceName string) (*trace.TracerProvider, error) {
	exporter, err := otlptracegrpc.New(context.Background(),
		otlptracegrpc.WithEndpoint(collectorEndpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	return trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(newResource(serviceName)),
	), nil
}

func newOTLPMeterProvider(collectorEndpoint, serviceName string) (*metric.MeterProvider, error) {
	exporter, err := otlpmetricgrpc.New(context.Background(),
		otlpmetricgrpc.WithEndpoint(collectorEndpoint),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	return metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(exporter)),
		metric.WithResource(newResource(serviceName)),
	), nil
}

func newOTLPLoggerProvider(collectorEndpoint, serviceName string) (*log.LoggerProvider, error) {
	exporter, err := otlploggrpc.New(context.Background(),
		otlploggrpc.WithEndpoint(collectorEndpoint),
		otlploggrpc.WithInsecure(),
	)

	if err != nil {
		return nil, err
	}

	return log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(exporter)),
		log.WithResource(newResource(serviceName)),
	), nil
}

func newResource(serviceName string) *resource.Resource {
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		// TODO: error out here this shouldn't be ignored silently
		return resource.Default()
	}

	return res
}
