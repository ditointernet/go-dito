package pubsub

import (
	"context"

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
