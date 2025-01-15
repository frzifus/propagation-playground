package main

import (
	"context"
	"fmt"
	"math/rand/v2"
	"net/http"
	"sync"
	"time"

	sdklogging "github.com/openshift-online/ocm-sdk-go/logging"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/trace"

	"github.com/frzifus/propagation-playground/internal/instr"
)

var tracer trace.Tracer

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

	otel.GetTracerProvider()
	tracer = otel.GetTracerProvider().Tracer("github.com/frzifus/propagation-playground/cmd/service2")
	db := &mockdb{}
	srv := http.Server{
		Addr:    "localhost:8080",
		Handler: otelhttp.NewMiddleware("incoming request")(endpoint(logger, db)),
	}
	go worker(logger, db)
	logger.Info(ctx, "listen and serve, addr: %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}

func endpoint(logger sdklogging.Logger, db database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		b := baggage.FromContext(ctx)
		bc := b.Member("correlationID").Value()
		br := b.Member("requestID").Value()
		logger.Info(ctx, "new incoming request, correlationID: %s, requestID: %s", bc, br)
		ctx, span := tracer.Start(ctx, "endpoint processing")
		defer span.End()
		// NOTE: do something
		time.Sleep(time.Duration((rand.Int32N(1501) + 500)) * time.Millisecond)
		j := job{
			name:     fmt.Sprintf("%s::%s::%s", time.Now().String(), bc, br),
			complete: false,
			modCtx:   span,
			bag:      b.Members(),
		}
		logger.Info(ctx, "create job, correlationID: %s, requestID: %s", bc, br)
		db.applyJob(ctx, j)
	}
}

type job struct {
	name     string
	complete bool

	// NOTE: other options?
	modCtx trace.Span
	bag    []baggage.Member
}

type database interface {
	// NOTE: mix create and update :D
	applyJob(ctx context.Context, j job)
	jobs(ctx context.Context) []job
}

type mockdb struct {
	mu      sync.Mutex
	interal []job
}

func (m *mockdb) applyJob(_ context.Context, j job) {
	m.mu.Lock()
	defer m.mu.Unlock()
	// Update...
	for i := 0; i < len(m.interal); i++ {
		if m.interal[i].name == j.name {
			m.interal[i] = j
			return
		}
	}
	// Not found, create..
	m.interal = append(m.interal, j)
}
func (m *mockdb) jobs(_ context.Context) []job {
	m.mu.Lock()
	defer m.mu.Unlock()
	jobs := []job{}
	for _, job := range m.interal {
		if !job.complete {
			jobs = append(jobs, job)
		}
	}
	return jobs
}

func worker(logger sdklogging.Logger, db database) {
	logger.Info(context.Background(), "worker: start")
	for {
		time.Sleep(5 * time.Second)
		jobs := db.jobs(context.Background())
		logger.Info(context.Background(), "worker: got  %d jobs", len(jobs))
		for _, job := range jobs {
			func() {
				ctx, span := tracer.Start(context.Background(), "worker-run")
				defer span.End()

				// NOTE: Create new baggag
				b, err := baggage.New(job.bag...)
				if err != nil {
					logger.Error(ctx, "failed to create baggage, err: %w", err)
					return
				}
				bc := b.Member("correlationID").Value()
				br := b.Member("requestID").Value()

				// NOTE: Create link
				attrs := []attribute.KeyValue{attribute.String("correlationID", bc), attribute.String("requestID", br)}
				link := trace.Link{
					SpanContext: job.modCtx.SpanContext(),
					Attributes:  attrs,
				}
				span.AddLink(link)
				span.SetAttributes(attrs...)
				logger.Info(ctx, "Execute Job: %s, correlationID: %s, requestID: %s", job.name, bc, br)
				job.complete = true
				db.applyJob(ctx, job)
				logger.Info(ctx, "Finished Job: %s, correlationID: %s, requestID: %s", job.name, bc, br)
			}()
		}
	}
}
