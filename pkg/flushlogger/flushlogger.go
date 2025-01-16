package flushlogger

import (
	"context"
	"io"
	"time"
)

type FlushLogger struct {
	writer         io.WriteCloser
	buf            chan []byte
	waitingForStop bool
	stopped        chan struct{}
}

func New(writer io.WriteCloser) *FlushLogger {
	logger := &FlushLogger{
		writer:         writer,
		buf:            make(chan []byte),
		waitingForStop: false,
		stopped:        make(chan struct{}),
	}

	go func() {
		for {
			if logger.waitingForStop && len(logger.buf) == 0 {
				close(logger.stopped)
				return
			}

			select {
			case v := <-logger.buf:
				//nolint:errcheck
				logger.writer.Write(v)
				//nolint:errcheck
				logger.writer.Write([]byte{'\n'})

			case <-time.After(20 * time.Millisecond):
			}
		}
	}()

	return logger
}

func (logger *FlushLogger) Write(v []byte) {
	logger.buf <- v
}

// Close closes the writer. Any accepted writes will be flushed. Any new writes will be rejected.
func (logger *FlushLogger) Close(ctx context.Context) {
	close(logger.buf)
	logger.waitingForStop = true

	select {
	case <-ctx.Done():
	case <-logger.stopped:
	}

	logger.writer.Close()
}
