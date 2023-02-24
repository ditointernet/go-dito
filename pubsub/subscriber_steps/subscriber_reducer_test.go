package steps_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"context"

	"github.com/golang/mock/gomock"

	"github.com/ditointernet/go-dito/errors"
	steps "github.com/ditointernet/go-dito/pubsub/subscriber_steps"
)

var _ = Describe("Reducer", func() {
	var (
		ctrl *gomock.Controller

		reduceFn     steps.ReduceFn
		initialState func() interface{}

		reducer steps.Reducer
	)

	BeforeEach(func() {
		t := GinkgoT()
		ctrl = gomock.NewController(t)

		reducer = steps.Reducer{}
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
			reducer.ReduceFn = reduceFn
			reducer.InitialState = initialState
			outCh = reducer.Do(ctx, inCh, errCh)
		})

		When("no input is being injected into pipeline", func() {
			It("should not produce any output", func() {
				Consistently(outCh).ShouldNot(Receive())
			})

			It("should not produce any error", func() {
				Consistently(errCh).ShouldNot(Receive())
			})
		})

		When("input channel contains lists of integer values", func() {
			BeforeEach(func() {
				go func() {
					inCh <- []int{0, 1, 2}
					inCh <- []int{3, 4, 5}
					inCh <- []int{6, 7, 8}
				}()
			})

			When("poorly manipulating input data", func() {
				var expectedError = errors.New("could not cast item into int")

				BeforeEach(func() {
					reduceFn = func(s, i interface{}, _ int) (interface{}, error) {
						_, ok := i.(string)
						if !ok {
							return nil, expectedError
						}

						return nil, errors.New("unexpected error")
					}

					initialState = func() interface{} {
						return ""
					}
				})

				It("should write into error channel", func() {
					Eventually(errCh).Should(Receive(Equal(expectedError)))
				})

				It("should not produce any output", func() {
					Consistently(outCh).ShouldNot(Receive())
				})
			})

			When("summing values", func() {
				BeforeEach(func() {
					reduceFn = func(s, i interface{}, _ int) (interface{}, error) {
						state, ok := s.(int)
						if !ok {
							return nil, errors.New("could not cast state into int")
						}

						item, ok := i.(int)
						if !ok {
							return nil, errors.New("could not cast item into int")
						}

						return state + item, nil
					}

					initialState = func() interface{} {
						return 0
					}
				})

				It("should sum the values of each input list", func() {
					Eventually(outCh).Should(Receive(Equal(3)))
					Eventually(outCh).Should(Receive(Equal(12)))
					Eventually(outCh).Should(Receive(Equal(21)))
				})

				It("should not produce any error", func() {
					Consistently(errCh).ShouldNot(Receive())
				})
			})
		})

		When("input channel of non-lists", func() {
			BeforeEach(func() {
				go func() {
					inCh <- 1
					inCh <- 2
					inCh <- 3
				}()
			})

			It("should write an ErrNonListValue into error channel", func() {
				Eventually(errCh).Should(Receive(Equal(steps.ErrNonListValue)))
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
