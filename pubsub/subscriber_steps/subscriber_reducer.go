package steps

import (
	"context"
	"errors"
	"reflect"
)

// ErrNonListValue is returned when a value that is not a list flows through the pipe at Reduce step.
var ErrNonListValue = errors.New("cannot map data that isn't an Slice")

// ReduceFn is the function that aggregates the data that passes through the pipeline into one final state.
type ReduceFn func(state interface{}, item interface{}, idx int) (newState interface{}, err error)

// Reducer is a Pubsub's subscriber pipeline step that aggregates a list of incoming pipe records into one.
type Reducer struct {
	ReduceFn     ReduceFn
	InitialState func() interface{}

	state interface{}
}

// Do executes a Reduce pipeline.
func (s Reducer) Do(ctx context.Context, inCh chan interface{}, errCh chan error) chan interface{} {

	chOut := make(chan interface{})
	go func() {
		defer close(chOut)

		for {
			select {
			case <-ctx.Done():
				return
			case in, ok := <-inCh:
				if !ok {
					return
				}

				out, err := s.do(ctx, in)
				if err != nil {
					errCh <- err
					break
				}

				chOut <- out
			}
		}
	}()

	return chOut
}

func (s *Reducer) do(ctx context.Context, in interface{}) (interface{}, error) {
	s.state = s.InitialState()

	switch reflect.TypeOf(in).Kind() {
	case reflect.Slice:
		slice := reflect.ValueOf(in)

		for idx := 0; idx < slice.Len(); idx++ {
			item := slice.Index(idx).Interface()

			newState, err := s.ReduceFn(s.state, item, idx)
			if err != nil {
				return nil, err
			}

			s.state = newState
		}

		return s.state, nil
	default:
		return nil, ErrNonListValue
	}
}
