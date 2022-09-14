package steps_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"context"

	"cloud.google.com/go/pubsub"
	"github.com/golang/mock/gomock"

	"github.com/ditointernet/go-dito/pubsub/mocks"
	steps "github.com/ditointernet/go-dito/pubsub/subscriber_steps"
)

var _ = Describe("MessageReceiver", func() {
	var (
		ctrl *gomock.Controller

		subscription *mocks.MockPubsubSubscriber

		receiver steps.SubscriberReceiver
	)

	BeforeEach(func() {
		t := GinkgoT()
		ctrl = gomock.NewController(t)

		subscription = mocks.NewMockPubsubSubscriber(ctrl)

		receiver = steps.SubscriberReceiver{
			Subscription: subscription,
		}
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("Do", func() {
		var (
			ctx context.Context

			inCh  chan interface{}
			errCh chan error
			outCh chan interface{}
		)

		BeforeEach(func() {
			ctx = context.Background()
			inCh = make(chan interface{})
			errCh = make(chan error)
		})

		JustBeforeEach(func() {
			outCh = receiver.Do(ctx, inCh, errCh)
		})

		Context("receiving messages from Pubsub", func() {
			When("an unexpected error occurs", func() {
				BeforeEach(func() {
					subscription.EXPECT().Receive(ctx, gomock.Any()).Return(ErrMock).Times(1)
				})

				It("should write an error into error channel", func() {
					Eventually(errCh).Should(Receive(&ErrMock))
				})

				It("should not write anything into output channel", func() {
					Consistently(outCh).ShouldNot(Receive())
				})

				It("should close output channel", func() {
					<-errCh
					Eventually(outCh).Should(BeClosed())
				})
			})

			When("receiving messages successfully", func() {
				mockedMessage := &pubsub.Message{
					Data: []byte("mocked message"),
				}

				BeforeEach(func() {
					subscription.EXPECT().Receive(ctx, gomock.Any()).
						DoAndReturn(func(ctx context.Context, callback func(context.Context, *pubsub.Message)) error {
							callback(ctx, mockedMessage)
							return nil
						}).
						Times(1)
				})

				It("should write a message into output channel", func() {
					Eventually(outCh).Should(Receive(Equal(mockedMessage)))
				})

				It("should not write anything into error channel", func() {
					Consistently(errCh).ShouldNot(Receive())
				})

				It("should close output channel", func() {
					Eventually(outCh).Should(BeClosed())
				})
			})
		})
	})
})
