package flushworker

import (
	"context"
)

type FlushWorker[T any] struct {
	buf     chan T
	handle  func(v T)
	stopped chan struct{}
}

func New[T any](handle func(v T)) *FlushWorker[T] {
	worker := &FlushWorker[T]{
		buf:     make(chan T),
		handle:  handle,
		stopped: make(chan struct{}),
	}

	go func() {
		for v := range worker.buf {
			worker.handle(v)
		}
		close(worker.stopped)
	}()

	return worker
}

func (worker *FlushWorker[T]) AddJob(v T) {
	worker.buf <- v
}

// Close closes the write to the buffer. Any accepted writes will be flushed. Any new writes will be rejected.
func (worker *FlushWorker[T]) Close(ctx context.Context) {
	close(worker.buf)

	select {
	case <-ctx.Done():
	case <-worker.stopped:
	}
}
