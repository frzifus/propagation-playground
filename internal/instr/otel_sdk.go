package instr

import (
	"context"
	"fmt"
	"time"

	sdklogging "github.com/openshift-online/ocm-sdk-go/logging"
	"go.opentelemetry.io/contrib/exporters/autoexport"
	"go.opentelemetry.io/otel/propagation"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"

	// semconv "go.opentelemetry.io/otel/semconv/v1.25.0"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
)

func InstallOpenTelemetryTracer(ctx context.Context, logger sdklogging.Logger) (func(context.Context) error, error) {
	logger.Info(ctx, "initialising OpenTelemetry tracer")

	exp, err := autoexport.NewSpanExporter(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTEL exporter: %w", err)
	}

	resources, err := resource.New(context.Background(),
		resource.WithAttributes(
		// semconv.ServiceNameKey.String("info.APPName"),
		// semconv.ServiceVersionKey.String("info.Version"),
		),
		resource.WithHost(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialise trace resources: %w", err)
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resources),
	)
	otel.SetTracerProvider(tp)

	shutdown := func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		return tp.Shutdown(ctx)
	}

	propagator := propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{})
	otel.SetTextMapPropagator(propagator)

	otel.SetErrorHandler(otelErrorHandlerFunc(func(err error) {
		logger.Error(ctx, "OpenTelemetry.ErrorHandler: %v", err)
	}))

	return shutdown, nil
}

type otelErrorHandlerFunc func(error)

// Handle implements otel.ErrorHandler
func (f otelErrorHandlerFunc) Handle(err error) {
	f(err)
}
