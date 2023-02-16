package trace

import (
	"context"

	"github.com/ditointernet/go-dito/errors"

	gcpexporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	otrace "go.opentelemetry.io/otel/trace"
)

// Params encodes necessary input data to initialize a new Tracer.
type Params struct {
	IsProductionEnvironment bool
	ApplicationName         string

	// ProjectID is an optional parameter, used for specifying which project the user is uploading the stats data to.
	// If ommited, this will default to your "Application Default Credentials".
	ProjectID string

	// TraceRatio indicates how often the system should collect traces.
	// Use it with caution: It may overload the system and also be too expensive to mantain its value too high in a
	// high throuput system. Values vary between 0 and 1, with 0 meaning No Sampling and 1 meaning Always Sampling.
	// Values lower than 0 are treated as 0 and values greater than 1 are treated as 1.
	TraceRatio float64
}

// NewTracer creates a new Tracer.
// It produces the tracer it self and a flush function that should be used to deliver any trace residue in cases of
// system shutdown. If your application is running outside of Google Cloud, make sure that your
// `GOOGLE_APPLICATION_CREDENTIALS` env variable is properly set.
func NewTracer(params Params) (otrace.Tracer, func(context.Context) error, error) {
	if params.ApplicationName == "" {
		return nil, nil, errors.NewMissingRequiredDependency("ApplicationName")
	}

	tOpts := []sdktrace.TracerProviderOption{
		sdktrace.WithResource(resource.NewSchemaless(
			attribute.KeyValue{
				Key:   semconv.ServiceNameKey,
				Value: attribute.StringValue(params.ApplicationName),
			},
		)),
	}

	if params.IsProductionEnvironment {
		exporter, err := gcpexporter.New(
			// Defaults to application credential's ProjectID if param is empty.
			gcpexporter.WithProjectID(params.ProjectID),
		)
		if err != nil {
			return nil, nil, err
		}

		tOpts = append(tOpts, sdktrace.WithSampler(sdktrace.TraceIDRatioBased(params.TraceRatio)))
		tOpts = append(tOpts, sdktrace.WithBatcher(exporter))
	}

	tp := sdktrace.NewTracerProvider(tOpts...)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tp.Tracer(params.ApplicationName), tp.Shutdown, nil
}

// MustNewTracer creates a new Tracer.
// It produces the tracer it self and a flush function that should be used to deliver any trace residue in cases of
// system shutdown. It panics if any error is found during tracer construction. If your application is running outside
// of Google Cloud, make sure that your `GOOGLE_APPLICATION_CREDENTIALS` env variable is properly set.
func MustNewTracer(params Params) (otrace.Tracer, func(context.Context) error) {
	tracer, flush, err := NewTracer(params)
	if err != nil {
		panic(err)
	}

	return tracer, flush
}
