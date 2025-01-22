package workerpool

import (
	"sync"
)

type WorkerPool[T any] struct {
	buf     chan T
	handle  func(v T, goRoutineId int)
	stopped chan struct{}
}

func New[T any](handle func(v T, goRoutineId int), concurrency int, bufferLength int) *WorkerPool[T] {
	workerPool := &WorkerPool[T]{
		buf:     make(chan T, bufferLength),
		handle:  handle,
		stopped: make(chan struct{}),
	}
	wg := new(sync.WaitGroup)
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func(goRoutineId int) {
			for v := range workerPool.buf {
				workerPool.handle(v, goRoutineId)
			}
			wg.Done()
		}(i)
	}

	go func() {
		wg.Wait()
		close(workerPool.stopped)
	}()

	return workerPool
}

func (workerPool *WorkerPool[T]) Enqueue(v T) {
	workerPool.buf <- v
}

// Close closes the write to the buffer. Any accepted writes will be flushed. Any new writes will be rejected.
func (workerPool *WorkerPool[T]) Close() {
	close(workerPool.buf)
	<-workerPool.stopped
}
