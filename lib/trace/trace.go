package trace

import (
	"context"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	otrace "go.opentelemetry.io/otel/trace"

	"github.com/ditointernet/go-dito/lib/errors"
)

// Params encondes necessary input data to initialize a new Tracer.
type Params struct {
	IsProductionEnvironment bool

	// ApplicationName should be in the following format: github.com/<project_name>/<repository>
	// E.G.: github.com/ditointernet/new-segments-service
	ApplicationName string

	// TraceRatio indicates how often the system should collect traces.
	// Use it with caution: It may overload the system and also be too expensive to mantain its value too high in a high throuput system
	// Values vary between 0 and 1, with 0 meaning No Sampling and 1 meaning Always Sampling.
	// Values lower than 0 are treated as 0 and values greater than 1 are treated as 1.
	TraceRatio float64
}

// NewTracer creates a new Tracer.
// It produces the tracer it self and a flush function that should be used to deliver any trace residue in cases of system shutdown.
// If your application is running outside of Google Cloud, make sure that your `GOOGLE_APPLICATION_CREDENTIALS` env variable is properly set.
func NewTracer(params Params) (otrace.Tracer, func(context.Context) error, error) {
	if params.ApplicationName == "" {
		return nil, nil, errors.NewMissingRequiredDependency("ApplicationName")
	}

	if !params.IsProductionEnvironment {
		tracer := otrace.NewNoopTracerProvider().Tracer(params.ApplicationName)
		return tracer, func(context.Context) error { return nil }, nil
	}

	exporter, err := texporter.NewExporter()
	if err != nil {
		return nil, nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(params.TraceRatio)),
		sdktrace.WithSyncer(exporter),
	)
	otel.SetTracerProvider(tp)

	return otel.GetTracerProvider().Tracer(params.ApplicationName), exporter.Shutdown, nil
}

// NewTracer creates a new Tracer.
// It produces the tracer it self and a flush function that should be used to deliver any trace residue in cases of system shutdown.
// It panics if any error is found during tracer construction.
// If your application is running outside of Google Cloud, make sure that your `GOOGLE_APPLICATION_CREDENTIALS` env variable is properly set.
func MustNewTracer(params Params) (otrace.Tracer, func(context.Context) error) {
	tracer, flush, err := NewTracer(params)
	if err != nil {
		panic(err)
	}

	return tracer, flush
}
