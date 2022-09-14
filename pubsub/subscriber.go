package pubsub

import (
	"context"

	"cloud.google.com/go/pubsub"
	"github.com/ditointernet/go-dito/errors"
	steps "github.com/ditointernet/go-dito/pubsub/subscriber_steps"
)

// DefaultReceiveSettings defines default settings for subscription ReceiveSettings.
type DefaultReceiveSettings struct {
	maxOutstandingMessages int
	Synchronous            bool
}

// DefaultReceiveSettings holds the default values for DefaultReceiveSettings.
var DefaultReceiveSettingsValues = DefaultReceiveSettings{
	maxOutstandingMessages: 500,
	Synchronous:            true,
}

// PubsubSubscriberPipelineParams encapsulates dependencies for a SubscriberPipelineParams instance.
type SubscriberPipelineParams struct {
	// PubsubClient is a Pubsub client, which interacts with Pubsub resources.
	PubsubClient *pubsub.Client

	// SubID is the unique identifier of the subscription within its project.
	SubID string

	// MaxOutstandingMessages is the maximum number of unprocessed messages
	// (unacknowledged but not yet expired). If MaxOutstandingMessages is 0,
	// the default value of DefaultReceiveSettingsValues.maxOutstandingMessages will be applied.
	// If the value is negative, then there will be no limit on the number of
	// unprocessed messages.
	MaxOutstandingMessages int
}

// SubscriberPipeline is a structure that defines a pubsub pipeline data handler.
type SubscriberPipeline interface {
	// Run executes the pipeline, connecting each registered step in a ordered way.
	Run(ctx context.Context) chan any

	// Map registers a new Mapper step into pipeline, which is modifies the data that passes
	// through the pipeline. It panics if any required dependency is not properly given.
	Map(mapFn steps.MapFn) SubscriberPipeline
}

// PipelineStep indicates how pipeline steps should execute each interaction with the pipe.
type PipelineStep interface {
	// Do executes a pipe entry.
	Do(context.Context, chan any, chan error) chan any
}

type subscriberPipeline struct {
	errCh chan error

	steps []PipelineStep
}

// NewPubsubSubscriberPipeline creates a new instance of pubsubSubscriberPipeline.
// The pipeline initiates with only one step: subscriberReceiver (which receives raw messages from Pubsub).
func NewSubscriberPipeline(params SubscriberPipelineParams) (subscriberPipeline, error) {
	if params.PubsubClient == nil {
		return subscriberPipeline{}, errors.NewMissingRequiredDependency("PubsubClient")
	}

	if params.SubID == "" {
		return subscriberPipeline{}, errors.NewMissingRequiredDependency("SubID")
	}

	maxOutstandingMessages := DefaultReceiveSettingsValues.maxOutstandingMessages
	if params.MaxOutstandingMessages == 0 {
		maxOutstandingMessages = params.MaxOutstandingMessages
	}

	sub := params.PubsubClient.Subscription(params.SubID)
	sub.ReceiveSettings.MaxOutstandingMessages = maxOutstandingMessages
	sub.ReceiveSettings.Synchronous = DefaultReceiveSettingsValues.Synchronous

	firstStep := steps.SubscriberReceiver{
		Subscription: sub,
	}

	sp := subscriberPipeline{
		errCh: make(chan error),
		steps: []PipelineStep{
			firstStep,
		},
	}

	return sp, nil
}

// MustNewSubscriberPipeline initializes subscriberPipeline by calling NewSubscriberPipeline
// It panics if any error is found.
func MustNewSubscriberPipeline(params SubscriberPipelineParams) subscriberPipeline {
	pb, err := NewSubscriberPipeline(params)
	if err != nil {
		panic(err)
	}

	return pb
}

// Run kicks off all pipeline steps executions, starting the subscription message receiving
// process and connecting each registered step in a ordered way.
func (sp subscriberPipeline) Run(ctx context.Context) chan any {
	// Spins up the receiver, which retrieves raw pubsub messages.
	ch := sp.steps[0].Do(ctx, nil, sp.errCh)

	// Attaches all additional steps to the pipeline.
	for ii := 1; ii < len(sp.steps); ii++ {
		ch = sp.steps[ii].Do(ctx, ch, sp.errCh)
	}

	// Fully configured channel, of which messages go through all pipeline steps.
	return ch
}

// Map registers a new Mapper step into pipeline.
// It panics if any required dependency is not properly given.
func (sp subscriberPipeline) Map(mapFn steps.MapFn) SubscriberPipeline {
	if mapFn == nil {
		panic(errors.NewMissingRequiredDependency("MapFn"))
	}

	sp.steps = append(sp.steps, steps.Mapper{MapFn: mapFn})
	return sp
}
