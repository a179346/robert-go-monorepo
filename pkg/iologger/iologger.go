package iologger

import (
	"encoding/json"
	"io"

	"github.com/a179346/robert-go-monorepo/pkg/gohf_extended"
	"github.com/a179346/robert-go-monorepo/pkg/workerpool"
)

type IoLogger struct {
	workerPool *workerpool.WorkerPool[gohf_extended.ApiLogData]
	writer     io.WriteCloser
}

func New(writer io.WriteCloser) *IoLogger {
	workerPool := workerpool.New(func(logData gohf_extended.ApiLogData, goRoutineId int) {
		bytes, err := json.Marshal(logData)
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

func (logger *IoLogger) Dispatch(logData gohf_extended.ApiLogData) {
	logger.workerPool.Enqueue(logData)
}

func (logger *IoLogger) Close() {
	logger.workerPool.Close()
	logger.writer.Close()
}
