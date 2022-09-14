package steps

import (
	"context"

	"cloud.google.com/go/pubsub"
)

type SubscriberReceiver struct {
	Receiver messageReceiver
}

type messageReceiver interface {
	// Receive calls f with the outstanding messages from the subscription.
	// It blocks until ctx is done, or the service returns a non-retryable error.
	Receive(ctx context.Context, f func(context.Context, *pubsub.Message)) error
}

// Do executes the messageReceiver pipeline step.
func (sr SubscriberReceiver) Do(ctx context.Context, _ chan interface{}, errCh chan error) chan interface{} {
	msgsCh := make(chan interface{})

	go func() {
		err := sr.Receiver.Receive(ctx, func(c context.Context, msg *pubsub.Message) {
			msgsCh <- msg
		})
		if err != nil {
			errCh <- err
		}

		close(msgsCh)
	}()

	return msgsCh
}
