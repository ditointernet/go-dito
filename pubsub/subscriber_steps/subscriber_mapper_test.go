package steps_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"context"

	"github.com/golang/mock/gomock"

	"github.com/ditointernet/go-dito/errors"

	steps "github.com/ditointernet/go-dito/pubsub/subscriber_steps"
)

var _ = Describe("Mapper", func() {
	var (
		ctrl *gomock.Controller

		mapFn steps.MapFn

		mapper steps.Mapper
	)

	BeforeEach(func() {
		t := GinkgoT()
		ctrl = gomock.NewController(t)

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
			mapper.MapFn = mapFn
			outCh = mapper.Do(ctx, inCh, errCh)
		})

		When("no input is being injected into pipeline", func() {
			It("should not produce any output", func() {
				Consistently(outCh).ShouldNot(Receive())
			})

			It("should not produce any error", func() {
				Consistently(errCh).ShouldNot(Receive())
			})
		})

		When("input channel with integer values", func() {
			const numItems = 10

			BeforeEach(func() {
				fillIntegerChannel(inCh, numItems)
			})

			When("doubleing input value", func() {
				BeforeEach(func() {
					mapFn = func(i interface{}) (interface{}, error) {
						v, ok := i.(int)
						if !ok {
							return nil, errors.New("could not cast input value into an integer")
						}

						return v * 2, nil
					}
				})

				It("should double each input value", func() {
					for i := 0; i < numItems; i++ {
						Eventually(outCh).Should(Receive(Equal(i * 2)))
					}
				})

				It("should not produce any error", func() {
					Consistently(errCh).ShouldNot(Receive())
				})
			})
		})

		When("poorly manipulating input data", func() {
			var expectedError = errors.New("could not cast input value into an string")

			BeforeEach(func() {
				fillIntegerChannel(inCh, 10)

				mapFn = func(i interface{}) (interface{}, error) {
					v, ok := i.(string)
					if !ok {
						return nil, expectedError
					}

					return v, nil
				}
			})

			It("should write into error channel", func() {
				Eventually(errCh).Should(Receive(Equal(expectedError)))
			})

			It("should not produce any output", func() {
				Consistently(outCh).ShouldNot(Receive())
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

func fillIntegerChannel(ch chan any, numItems int) {
	go func() {
		for i := 0; i < numItems; i++ {
			ch <- i
		}
	}()
}
