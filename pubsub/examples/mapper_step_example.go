package examples

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	godito "github.com/ditointernet/go-dito/pubsub"
)

// Mapper_step_example exhibits the operation of a recently instantiated SubscriberPipeline with
// an additional mapper step attached, which modifies the data that goes through the pipeline.
//
// The pipeline outputs the contents of a channel, which is channeling *Pubsub.Message messages,
// in this case. The message then goes through a Mapper step, which is suported by the package,
// but provided by the client.
//
// A Mapper is particularly useful for consuming the raw *Pubsub.Message contents and transforming
// them into a more practical custom type (such as myCustomType).
func Mapper_step_example() {
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

	pipelineWithMapper := pipeline.Map(myMapper)

	for in := range pipelineWithMapper.Run(ctx) {
		msg, ok := in.(myCustomType)
		if !ok {
			fmt.Println("myCustomType type casting wasn't successful")
			return
		}

		fmt.Println(msg.myData)
	}
}

type myCustomType struct {
	myData string
}

func myMapper(in any) (any, error) {
	msg, ok := in.(*pubsub.Message)
	if !ok {
		return nil, fmt.Errorf("*pubsub.Message type casting wasn't successful")
	}

	out := myCustomType{
		myData: fmt.Sprintf("my modified data: %s", string(msg.Data)),
	}

	return out, nil
}
