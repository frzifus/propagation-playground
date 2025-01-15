package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	sdklogging "github.com/openshift-online/ocm-sdk-go/logging"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"

	"github.com/frzifus/propagation-playground/internal/instr"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logger, err := sdklogging.NewStdLoggerBuilder().Build()
	if err != nil {
		panic(err)
	}
	shutdown, err := instr.InstallOpenTelemetryTracer(ctx, logger)
	defer shutdown(ctx)
	if err != nil {
		logger.Fatal(ctx, "could not install opentelemetry tracer, err: %w", err)
	}

	client := &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	tracer := otel.GetTracerProvider().Tracer("github.com/frzifus/propagation-playground/cmd/service1")

	correlationID := uuid.New()
	for requestID := 0; requestID < 10; requestID++ {
		cID, _ := baggage.NewMember("correlationID", correlationID.String())
		rID, _ := baggage.NewMember("requestID", fmt.Sprintf("%d", requestID))

		bag, err := baggage.New(cID, rID)
		if err != nil {
			logger.Error(ctx, "failed to create baggage, err: %w", err)
			continue
		}

		func(ctx context.Context) {
			ctx, span := tracer.Start(ctx, "origin")
			defer span.End()
			b := baggage.FromContext(ctx)
			bc := b.Member("correlationID").Value()
			br := b.Member("requestID").Value()
			span.SetAttributes(attribute.String("correlationID", bc))
			span.SetAttributes(attribute.String("requestID", br))
			r, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080", nil)
			if err != nil {
				logger.Error(ctx, "could not creat request, correlationID: %s, requestID: %s, err: %v", bc, br, err)
				return
			}
			logger.Info(ctx, "Do request, correlationID: %s, requestID: %s", bc, br)
			resp, err := client.Do(r)
			if err != nil {
				logger.Error(ctx, "something went wrong, correlationID: %s, requestID: %s, err: %v", bc, br, err)
				return
			}
			logger.Info(ctx, "status: %d, correlationID: %s, requestID: %s", resp.StatusCode, bc, br)
		}(baggage.ContextWithBaggage(ctx, bag))
		time.Sleep(500 * time.Millisecond)
	}
}
