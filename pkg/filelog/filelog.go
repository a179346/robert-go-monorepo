package filelog

import (
	"context"
	"encoding/json"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

type FileLog struct {
	writer         *lumberjack.Logger
	buf            chan interface{}
	waitingForStop bool
	stopped        chan struct{}
}

func New(logger *lumberjack.Logger) *FileLog {
	log := &FileLog{
		writer:         logger,
		buf:            make(chan interface{}),
		waitingForStop: false,
		stopped:        make(chan struct{}),
	}

	go func() {
		for {
			if log.waitingForStop && len(log.buf) == 0 {
				close(log.stopped)
				return
			}

			select {
			case v := <-log.buf:
				//nolint:errcheck
				json.NewEncoder(log.writer).Encode(v)

			case <-time.After(20 * time.Millisecond):
			}
		}
	}()

	return log
}

func (log *FileLog) JSON(v interface{}) {
	log.buf <- v
}

// Close close the write to buffer and wait until buffer is empty
func (log *FileLog) Close(ctx context.Context) {
	close(log.buf)
	log.waitingForStop = true

	select {
	case <-ctx.Done():
	case <-log.stopped:
	}

	log.writer.Close()
}
