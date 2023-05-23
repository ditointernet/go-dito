package steps_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"context"
	"reflect"
	"time"

	"github.com/golang/mock/gomock"

	steps "github.com/ditointernet/go-dito/pubsub/subscriber_steps"
)

var _ = Describe("Batcher", func() {
	const (
		flushTimeout = time.Millisecond * 10
		batchSize    = 5
	)

	var (
		ctrl *gomock.Controller

		batcher steps.Batcher
	)

	BeforeEach(func() {
		t := GinkgoT()
		ctrl = gomock.NewController(t)

		batcher = steps.Batcher{
			BatchSize: batchSize,
			Timeout:   flushTimeout,
			ItemType:  reflect.TypeOf(1), // Batching integers
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
			outCh = batcher.Do(ctx, inCh, errCh)
		})

		When("no input is being injected into pipeline", func() {
			It("should not produce any output", func() {
				Consistently(outCh).ShouldNot(Receive())
			})

			It("should not produce any error", func() {
				Consistently(errCh).ShouldNot(Receive())
			})
		})

		When("inputing less data then the size of the batch", func() {
			BeforeEach(func() {
				fillIntegerChannel(inCh, batchSize-1)
			})

			It("should not produce any error", func() {
				Consistently(errCh).ShouldNot(Receive())
			})

			It("should produce only one batch with 4 (batchSize - 1) items", func() {
				Eventually(outCh).Should(Receive(Equal([]int{0, 1, 2, 3})))
			})
		})

		When("inputing more data then the size of the batch", func() {
			BeforeEach(func() {
				fillIntegerChannel(inCh, batchSize*2)
			})

			It("should not produce any error", func() {
				Consistently(errCh).ShouldNot(Receive())
			})

			It("should produce two batchs with 5 (batchSize) items", func() {
				Eventually(outCh).Should(Receive(Equal([]int{0, 1, 2, 3, 4})))
				Eventually(outCh).Should(Receive(Equal([]int{5, 6, 7, 8, 9})))
			})
		})

		When("input channel is closed", func() {
			BeforeEach(func() {
				close(inCh)
			})

			It("should close its output channel", func() {
				Eventually(outCh).Should(BeClosed())
			})
		})

		When("context is done", func() {
			BeforeEach(func() {
				var cancelFn context.CancelFunc
				ctx, cancelFn = context.WithCancel(ctx)

				cancelFn()
			})

			It("should close its output channel", func() {
				Eventually(outCh).Should(BeClosed())
			})
		})
	})
})
