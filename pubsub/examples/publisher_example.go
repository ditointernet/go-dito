package examples

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"cloud.google.com/go/pubsub"
	godito "github.com/ditointernet/go-dito/pubsub"
)

// MessageSchema represents the message schema that will be published.
type MessageSchema struct {
	Attr string
}

// ToBytes marshals itself using it's instance data.
func (ms MessageSchema) ToBytes() ([]byte, error) {
	data, err := json.Marshal(ms)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Publisher_pipeline_example shows the operation of a recently instantiated Generics compatible PubsubClient publisher.
// It accepts any message schema, requiring only that it's type implements ToByteser interface. The message schema,
// along with a wrapped Pubsub Topic (created with a NewTopicWrapper), must be passed to the PubsubClient builder
// (MustNewPubSubClient).
func Publisher_pipeline_example() {
	PROJECT_ID := "dito-it-tracking-dev"
	TOPIC_ID := "publisher_test"

	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, PROJECT_ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close()

	topic := client.Topic(TOPIC_ID)

	publisher := godito.MustNewPubSubClient[MessageSchema](godito.NewTopicWrapper(topic))

	var inputList []godito.PublishInput[MessageSchema]
	for i := 0; i < 2; i++ {
		in := godito.PublishInput[MessageSchema]{
			Data: MessageSchema{
				Attr: fmt.Sprintf("fake-publish-data-%s", strconv.Itoa(i)),
			},
			Attributes: map[string]string{
				"test": "test",
			},
		}

		inputList = append(inputList, in)
	}

	publisher.Publish(ctx, inputList...)
}
