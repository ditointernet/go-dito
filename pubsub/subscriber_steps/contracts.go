package steps

import (
	"context"

	"cloud.google.com/go/pubsub"
)

// PubsubSubscription defines something that knows how to receive Pubsub messages
// just like a Pubsub Subscription would.
type PubsubSubscription interface {
	// Receive calls f with the outstanding messages from the subscription.
	// It blocks until ctx is done, or the service returns a non-retryable error.
	Receive(ctx context.Context, f func(context.Context, *pubsub.Message)) error
}
