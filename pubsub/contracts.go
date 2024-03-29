package pubsub

import (
	"context"
	"reflect"
	"time"

	"cloud.google.com/go/pubsub"
)

// Publisher defines boundary interfaces of a pubsub topic.
type Publisher interface {
	Publish(ctx context.Context, msg *pubsub.Message) Getter
}

// Getter defines boundary interfaces of a pubsub result object.
type Getter interface {
	Get(ctx context.Context) (serverID string, err error)
}

// ToByteser defines the interface of pubsub client types.
type ToByteser interface {
	ToBytes() ([]byte, error)
}

// SubscriberPipeline is a structure that defines a pubsub pipeline data handler.
type SubscriberPipeline interface {
	// Run executes the pipeline, connecting each registered step in a ordered way.
	Run(ctx context.Context) chan any

	// Map registers a new Mapper step into pipeline, which is modifies the data that passes
	// through the pipeline. It panics if any required dependency is not properly given.
	Map(mapFn func(any) (any, error)) SubscriberPipeline

	// Reduce registers a new Reducer step into pipeline.
	// It panics if any required dependency is not properly given.
	Reduce(reduceFn func(state interface{}, item interface{}, idx int) (newState interface{}, err error), initialState func() interface{}) SubscriberPipeline

	// Batch registers a new Batcher step into pipeline.
	// It panics if any required dependency is not properly given.
	Batch(itemType reflect.Type, batchSize int, timeout time.Duration) SubscriberPipeline

	// Errors exposes all errors that happens during pipeline processing.
	Errors() chan error
}

// Doer indicates how pipeline steps should execute each interaction with the pipe.
type Doer interface {
	// Do executes a pipe entry.
	Do(context.Context, chan any, chan error) chan any
}
