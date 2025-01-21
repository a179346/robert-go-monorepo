package flushworker

import (
	"sync"
)

type FlushWorker[T any] struct {
	buf     chan T
	handle  func(v T, goRoutineId int)
	stopped chan struct{}
}

func New[T any](handle func(v T, goRoutineId int), concurrency int, bufferLength int) *FlushWorker[T] {
	worker := &FlushWorker[T]{
		buf:     make(chan T, bufferLength),
		handle:  handle,
		stopped: make(chan struct{}),
	}
	wg := new(sync.WaitGroup)
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			for v := range worker.buf {
				worker.handle(v, i)
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(worker.stopped)
	}()

	return worker
}

func (worker *FlushWorker[T]) AddJob(v T) {
	worker.buf <- v
}

// Close closes the write to the buffer. Any accepted writes will be flushed. Any new writes will be rejected.
func (worker *FlushWorker[T]) Close() {
	close(worker.buf)
	<-worker.stopped
}
