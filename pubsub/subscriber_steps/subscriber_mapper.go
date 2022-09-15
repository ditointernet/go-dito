package steps

import (
	"context"
)

// MapFn is the function that modifies the data that passes through the pipeline.
type MapFn func(any) (any, error)

// Mapper is a pipeline step that modifies each record that passes through the pipeline.
type Mapper struct {
	MapFn MapFn
}

// Do executes a Map pipeline.
func (m Mapper) Do(ctx context.Context, inCh chan any, errCh chan error) chan any {
	outCh := make(chan any)

	go func() {
		defer close(outCh)

		for {
			select {
			case <-ctx.Done():
				return
			case in, ok := <-inCh:
				if !ok {
					return
				}

				out, err := m.do(ctx, in)
				if err != nil {
					errCh <- err
					break
				}

				outCh <- out
			}
		}
	}()

	return outCh
}

func (m Mapper) do(ctx context.Context, in any) (any, error) {
	out, err := m.MapFn(in)
	if err != nil {
		return nil, err
	}

	return out, nil
}
