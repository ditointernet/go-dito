package pubsub

import (
	"context"

	"cloud.google.com/go/pubsub"
)

// TopicWrapper envelopes a pubsub topic type.
type TopicWrapper struct {
	topic *pubsub.Topic
}

// Publish envelopes a pubsub topic publish method.
// It returns a Geter.
func (tw TopicWrapper) Publish(ctx context.Context, msg *pubsub.Message) Getter {
	result := tw.topic.Publish(ctx, msg)
	return result
}

// NewTopicWrapper returns a new instance of TopicWrapper.
func NewTopicWrapper(topic *pubsub.Topic) TopicWrapper {
	return TopicWrapper{
		topic: topic,
	}
}
