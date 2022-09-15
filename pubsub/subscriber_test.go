package pubsub

import (
	"context"

	"cloud.google.com/go/pubsub"
	"github.com/ditointernet/go-dito/errors"
	steps "github.com/ditointernet/go-dito/pubsub/subscriber_steps"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Subscriber", func() {
	Context("NewSubscriberPipeline", func() {
		var (
			mockErrCh chan error
		)

		BeforeEach(func() {
			mockErrCh = make(chan error)
		})

		When("subscription dependency is missing", func() {
			It("returns the an MissingRequiredDependency error", func() {
				pipe, err := NewSubscriberPipeline(SubscriberPipelineParams{})

				Expect(pipe).To(Equal(subscriberPipeline{}))
				Expect(err).To(Equal(errors.NewMissingRequiredDependency("PubsubSubscription")))
			})
		})

		When("all dependencies are provided", func() {
			It("returns a working pipeline", func() {
				pipe, err := NewSubscriberPipeline(SubscriberPipelineParams{
					PubsubSubscription: fakeSub{},
					errCh:              mockErrCh,
				})

				Expect(pipe).To(Equal(subscriberPipeline{
					errCh: mockErrCh,
					steps: []PipelineStep{
						steps.SubscriberReceiver{
							Subscription: fakeSub{},
						},
					},
				}))
				Expect(err).To(BeNil())
			})
		})

	})

	Context("MustNewSubscriberPipeline", func() {
		var (
			mockErrCh chan error
		)

		BeforeEach(func() {
			mockErrCh = make(chan error)
		})

		When("missing subscription", func() {
			It("panics with a PubsubSubscription MissingRequiredDependency error", func() {
				Expect(func() {
					_ = MustNewSubscriberPipeline(SubscriberPipelineParams{})

				}).To(PanicWith(Equal(errors.NewMissingRequiredDependency("PubsubSubscription"))))
			})
		})

		When("missing subscription", func() {
			It("panics with a PubsubSubscription MissingRequiredDependency error", func() {
				Expect(func() {
					_ = MustNewSubscriberPipeline(SubscriberPipelineParams{})

				}).To(PanicWith(Equal(errors.NewMissingRequiredDependency("PubsubSubscription"))))
			})
		})

		When("all dependencies are provided", func() {
			It("successfully creates a subscriber pipeline with no panic", func() {
				Expect(func() {
					_ = MustNewSubscriberPipeline(SubscriberPipelineParams{
						PubsubSubscription: fakeSub{},
						errCh:              mockErrCh,
					})
				}).NotTo(Panic())
			})
		})
	})
})

type fakeSub struct{}

func (ft fakeSub) Receive(ctx context.Context, f func(context.Context, *pubsub.Message)) error {
	var err error

	return err
}
