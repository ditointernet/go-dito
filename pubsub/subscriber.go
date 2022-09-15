package pubsub

import (
	"context"

	"github.com/ditointernet/go-dito/errors"
	steps "github.com/ditointernet/go-dito/pubsub/subscriber_steps"
)

// SubscriberPipelineParams encapsulates dependencies for a SubscriberPipelineParams instance.
type SubscriberPipelineParams struct {
	PubsubSubscription steps.Receiver

	errCh chan error
}

type subscriberPipeline struct {
	errCh chan error

	steps []Doer
}

// NewSubscriberPipeline creates a new instance of subscriberPipeline.
// The pipeline initiates with only one step: subscriberReceiver (which receives raw messages from Pubsub).
func NewSubscriberPipeline(params SubscriberPipelineParams) (subscriberPipeline, error) {
	if params.PubsubSubscription == nil {
		return subscriberPipeline{}, errors.NewMissingRequiredDependency("PubsubSubscription")
	}

	if params.errCh == nil {
		params.errCh = make(chan error)
	}

	firstStep := steps.SubscriberReceiver{
		Subscription: params.PubsubSubscription,
	}

	sp := subscriberPipeline{
		errCh: params.errCh,
		steps: []Doer{
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
// process and connecting each additional registered step in a ordered way.
func (sp subscriberPipeline) Run(ctx context.Context) chan any {
	// Spins up the receiver, which retrieves raw pubsub messages.
	ch := sp.steps[0].Do(ctx, nil, sp.errCh)

	// Attaches all additional steps to the pipeline.
	for index := 1; index < len(sp.steps); index++ {
		ch = sp.steps[index].Do(ctx, ch, sp.errCh)
	}

	// Fully configured channel, with messages that go through all pipeline steps.
	return ch
}

// Map registers a new Mapper step into pipeline.
// It panics if any required dependency is not properly given.
func (sp subscriberPipeline) Map(mapFn func(any) (any, error)) SubscriberPipeline {
	if mapFn == nil {
		panic(errors.NewMissingRequiredDependency("MapFn"))
	}

	sp.steps = append(sp.steps, steps.Mapper{MapFn: mapFn})
	return sp
}
