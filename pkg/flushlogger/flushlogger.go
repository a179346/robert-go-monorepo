package flushlogger

import (
	"context"
	"io"

	"github.com/a179346/robert-go-monorepo/pkg/flushworker"
)

type FlushLogger struct {
	worker *flushworker.FlushWorker[[]byte]
	writer io.WriteCloser
}

func New(writer io.WriteCloser) *FlushLogger {
	worker := flushworker.New(func(v []byte, goRoutineId int) {
		//nolint:errcheck
		writer.Write(append(v, '\n'))
	}, 1, 1024)

	return &FlushLogger{worker, writer}
}

func (logger *FlushLogger) Write(v []byte) {
	logger.worker.AddJob(v)
}

func (logger *FlushLogger) Close(ctx context.Context) {
	logger.worker.Close(ctx)
	logger.writer.Close()
}
