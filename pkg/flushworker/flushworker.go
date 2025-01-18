package flushworker

import (
	"context"
	"time"
)

type FlushWorker[T any] struct {
	buf            chan T
	handle         func(v T)
	waitingForStop bool
	stopped        chan struct{}
}

func New[T any](handle func(v T)) *FlushWorker[T] {
	worker := &FlushWorker[T]{
		buf:            make(chan T),
		handle:         handle,
		waitingForStop: false,
		stopped:        make(chan struct{}),
	}

	go func() {
		for {
			if worker.waitingForStop && len(worker.buf) == 0 {
				close(worker.stopped)
				return
			}

			select {
			case v := <-worker.buf:
				worker.handle(v)

			case <-time.After(20 * time.Millisecond):
			}
		}
	}()

	return worker
}

func (worker *FlushWorker[T]) AddJob(v T) {
	worker.buf <- v
}

// Close closes the write to the buffer. Any accepted writes will be flushed. Any new writes will be rejected.
func (worker *FlushWorker[T]) Close(ctx context.Context) {
	close(worker.buf)
	worker.waitingForStop = true

	select {
	case <-ctx.Done():
	case <-worker.stopped:
	}
}
