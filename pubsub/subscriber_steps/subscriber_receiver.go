package steps

import (
	"context"

	"cloud.google.com/go/pubsub"
)

// SubscriberReceiver is the first step of the pipeline which is responsible for reading the message from pubsub subscription.
type SubscriberReceiver struct {
	Subscription Receiver
}

// Do executes the messageReceiver pipeline step.
func (sr SubscriberReceiver) Do(ctx context.Context, _ chan interface{}, errCh chan error) chan interface{} {
	msgsCh := make(chan interface{})

	go func() {
		err := sr.Subscription.Receive(ctx, func(c context.Context, msg *pubsub.Message) {
			msgsCh <- msg
		})
		if err != nil {
			errCh <- err
		}

		close(msgsCh)
	}()

	return msgsCh
}
