package examples

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	godito "github.com/ditointernet/go-dito/pubsub"
)

// Raw_pipeline_example shows the raw operation of a recently instantiated SubscriberPipeline with
// no additional steps attached.
//
// The SubscriberPipeline builder (MustNewSubscriberPipeline) requires a Pubsub Subscriber
// (or something that equally implements it's Receive functionality). The pipeline then
// outputs the contents of a channel, which in it's turn, channels *Pubsub.Message messages
// in this case.
func Raw_pipeline_example() {
	PROJECT_ID := "your-project"
	SUB_ID := "your-subscription"

	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, PROJECT_ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close()

	sub := client.Subscription(SUB_ID)

	pipeline := godito.MustNewSubscriberPipeline(godito.SubscriberPipelineParams{
		PubsubSubscription: sub,
	})

	for in := range pipeline.Run(ctx) {
		msg, ok := in.(*pubsub.Message)
		if !ok {
			fmt.Println("*pubsub.Message type casting wasn't successful")
			return
		}

		fmt.Println(string(msg.Data))
	}
}
