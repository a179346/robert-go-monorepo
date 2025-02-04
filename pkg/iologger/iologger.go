package iologger

import (
	"encoding/json"
	"io"

	"github.com/a179346/robert-go-monorepo/pkg/apilog"
	"github.com/a179346/robert-go-monorepo/pkg/workerpool"
)

type IoLogger struct {
	workerPool *workerpool.WorkerPool[apilog.Data]
	writer     io.WriteCloser
}

func New(writer io.WriteCloser) *IoLogger {
	workerPool := workerpool.New(func(data apilog.Data, goRoutineId int) {
		bytes, err := json.Marshal(data)
		if err != nil {
			return
		}

		//nolint:errcheck
		writer.Write(append(bytes, '\n'))
	}, 1, 1024, true)

	return &IoLogger{
		workerPool: workerPool,
		writer:     writer,
	}
}

func (logger *IoLogger) Dispatch(data apilog.Data) {
	logger.workerPool.Enqueue(data)
}

func (logger *IoLogger) Close() {
	logger.workerPool.Close()
	logger.writer.Close()
}
