package pubsub

import (
	"context"
	"encoding/json"

	"go.opentelemetry.io/otel/trace"

	"cloud.google.com/go/pubsub"
)

// TraceIDContextKey defines the trace id key in a context.
const TraceIDContextKey string = "trace_id"

// PubSubClient is responsible for managing a pubsub topic.
type PubSubClient[T json.Marshaler] struct {
	topic Publisher
}

// NewPubSubClient returns a new instance of PubSubClient.
func NewPubSubClient[T json.Marshaler](topic Publisher) (PubSubClient[T], error) {
	return PubSubClient[T]{
		topic: topic,
	}, nil
}

// MustNewPubSubClient initializes Publisher by calling NewPubSubClient
// It panics if any error is found.
func MustNewPubSubClient[T json.Marshaler](topic Publisher) PubSubClient[T] {
	p, err := NewPubSubClient[T](topic)
	if err != nil {
		panic(err)
	}

	return p
}

// PublishInput is the input for publishing data in a topic.
type PublishInput[T json.Marshaler] struct {
	Data       T
	Attributes map[string]string
}

// Publish publishes messages in a pubsub topic.
func (c PubSubClient[T]) Publish(ctx context.Context, in ...PublishInput[T]) []error {
	var errs []error

	traceID := getTraceID(trace.SpanFromContext(ctx))

	for _, message := range in {
		data, err := message.Data.MarshalJSON()
		if err != nil {
			errs = append(errs, err)
			continue
		}

		message.Attributes[TraceIDContextKey] = traceID
		pubSubMsg := &pubsub.Message{
			Data:       data,
			Attributes: message.Attributes,
		}

		result := c.topic.Publish(ctx, pubSubMsg)
		_, err = result.Get(ctx)

		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func getTraceID(span trace.Span) string {
	if !span.SpanContext().HasTraceID() {
		return ""
	}

	return span.SpanContext().TraceID().String()
}
