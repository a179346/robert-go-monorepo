package flushlogger

import (
	"encoding/json"
	"io"

	"github.com/a179346/robert-go-monorepo/pkg/flushworker"
	"github.com/a179346/robert-go-monorepo/pkg/gohf_extended"
)

type FlushLogger struct {
	worker *flushworker.FlushWorker[gohf_extended.LogData]
	writer io.WriteCloser
}

func New(writer io.WriteCloser) *FlushLogger {
	worker := flushworker.New(func(logData gohf_extended.LogData, goRoutineId int) {
		bytes, err := json.Marshal(logData)
		if err != nil {
			return
		}

		//nolint:errcheck
		writer.Write(append(bytes, '\n'))
	}, 1, 1024)

	return &FlushLogger{worker, writer}
}

func (logger *FlushLogger) Write(logData gohf_extended.LogData) {
	logger.worker.AddJob(logData)
}

func (logger *FlushLogger) Close() {
	logger.worker.Close()
	logger.writer.Close()
}
