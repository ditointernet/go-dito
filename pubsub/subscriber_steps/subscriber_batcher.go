package steps

import (
	"context"
	"reflect"
	"time"
)

// Batcher is a Pubsub's subscriber pipeline step that accumulates messages in batches.
// It is useful in situations where there is a bunch of unitary messages that should be
// grouped to reduce systems internal I/Os, improving its performance and scale capabilities.
type Batcher struct {
	BatchSize int
	Timeout   time.Duration
	ItemType  reflect.Type

	batchItems reflect.Value
}

// Do executes a Batch pipeline.
func (s Batcher) Do(ctx context.Context, inCh chan interface{}, errCh chan error) chan interface{} {

	s.batchItems = reflect.MakeSlice(reflect.SliceOf(s.ItemType), 0, s.BatchSize)
	outCh := make(chan interface{})

	go func() {
		for {
			defer close(outCh)

			select {
			case in, ok := <-inCh:
				if !ok {
					return
				}

				s.batchItems = reflect.Append(s.batchItems, reflect.ValueOf(in))
				if s.shouldFlush() {
					s.flush(ctx, outCh)
				}
			case <-time.NewTimer(s.Timeout).C:
				if s.batchItems.Len() > 0 {
					s.flush(ctx, outCh)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return outCh
}

func (s Batcher) shouldFlush() bool {
	return s.batchItems.Len() >= s.BatchSize
}

func (s *Batcher) flush(ctx context.Context, outCh chan interface{}) {
	outCh <- s.batchItems.Interface()
	s.batchItems = reflect.MakeSlice(reflect.SliceOf(s.ItemType), 0, s.BatchSize)
}
