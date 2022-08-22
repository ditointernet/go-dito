package pubsub

import (
	"context"

	"cloud.google.com/go/pubsub"
)

// Topicier defines boundary interfaces of a pubsub topic.
type Topicer interface {
	Publish(ctx context.Context, msg *pubsub.Message) Resultier
}

// Resultier defines boundary interfaces of a pubsub result object.
type Resultier interface {
	Get(ctx context.Context) (serverID string, err error)
}
