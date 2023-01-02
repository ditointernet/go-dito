package examples

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"cloud.google.com/go/pubsub"
	godito "github.com/ditointernet/go-dito/pubsub"
)

func Batcher_step_example() {
	PROJECT_ID := "your-project"
	SUB_ID := "your-subscription"
	BATCH_SIZE := 100
	BATCH_MAX_FLUSH_TIMEOUT := time.Millisecond * time.Duration(500)

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

	pipelineWithBatch := pipeline.
		Map(JobMapper).
		Batch(reflect.TypeOf(ValidationJob{}), BATCH_SIZE, BATCH_MAX_FLUSH_TIMEOUT).
		Reduce(JobGrouper, func() interface{} { return ValidationJobGroups{} })

	for in := range pipelineWithBatch.Run(ctx) {
		msg, ok := in.(ValidationJobGroups)
		if !ok {
			fmt.Println("ValidationJobGroups type casting wasn't successful")
			return
		}

		fmt.Println(msg)
	}
}

// ValidationJob represents a domain ValidationMessage at Pubsub.
type ValidationJob struct {
	Attribute1 string

	Data1 string
	Data2 string
}

type ValidationData struct {
	Data1 string `json:"data_1"`
	Data2 string `json:"data_2"`
}

// ValidationJobList is a list of validation jobs
type ValidationJobList []ValidationJob

// ValidationJobGroups groups validation jobs by a key
type ValidationJobGroups map[string]ValidationJobList

// PubsubMessage encapsulates pubsub messages
type PubsubMessage struct {
	Msg *pubsub.Message
}

// NewPubsubValidationJob creates a new ValidationJob instance.
func NewPubsubValidationJob(pm PubsubMessage) (ValidationJob, error) {
	var data ValidationData
	if err := json.Unmarshal(pm.Msg.Data, &data); err != nil {
		return ValidationJob{}, err
	}

	return ValidationJob{
		Attribute1: pm.Msg.Attributes["attribute_1"],

		Data1: data.Data1,
		Data2: data.Data2,
	}, nil
}

// JobMapper transforms a input into a pubsub validation job
func JobMapper(in interface{}) (interface{}, error) {
	data, ok := in.(*pubsub.Message)
	if !ok {
		return nil, errors.New("could not transform pipeline data into a *pubsub.Message")
	}

	job, err := NewPubsubValidationJob(PubsubMessage{Msg: data})
	if err != nil {
		return nil, err
	}

	return job, nil
}

// JobGrouper groups jobs
func JobGrouper(s, i interface{}, idx int) (interface{}, error) {
	state, ok := s.(ValidationJobGroups)
	if !ok {
		return nil, errors.New("cannot build reduce state into a (map[string]message.ValidationJobList instance")
	}

	item, ok := i.(ValidationJob)
	if !ok {
		return nil, errors.New("cannot build reduce item into a message.ValidationJob instance")
	}

	key := item.Attribute1
	if state[key] == nil {
		state[key] = ValidationJobList{}
	}

	state[key] = append(state[key], item)

	return state, nil
}
