package rabbitmqlogger

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/a179346/robert-go-monorepo/pkg/gohf_extended"
	"github.com/a179346/robert-go-monorepo/pkg/workerpool"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQLogger struct {
	workerPool *workerpool.WorkerPool[gohf_extended.LogData]
	conn       *amqp.Connection
	channels   []*amqp.Channel
}

func New(conn *amqp.Connection, exchange string) *RabbitMQLogger {
	const concurrency = 8
	const maxRetryCnt = 5

	logger := &RabbitMQLogger{
		conn:     conn,
		channels: make([]*amqp.Channel, concurrency),
	}

	workerPool := workerpool.New(func(logData gohf_extended.LogData, goRoutineId int) {
		body, err := json.Marshal(logData)
		if err != nil {
			return
		}

		for retryCnt := 1; retryCnt <= maxRetryCnt; retryCnt++ {
			ch, err := logger.getChannel(goRoutineId)
			if err == nil {
				err = ch.Publish(
					exchange,
					"",
					false,
					false,
					amqp.Publishing{
						MessageId:    logData.ID,
						AppId:        logData.App,
						ContentType:  "application/json",
						DeliveryMode: 2,
						Body:         body,
						Timestamp:    time.UnixMilli(logData.StartUnixMs),
					},
				)
			}

			if err == nil {
				return
			}
			time.Sleep(time.Duration(retryCnt*2) * time.Second)
		}
	}, concurrency, 1024)
	logger.workerPool = workerPool

	return logger
}

func (logger *RabbitMQLogger) Dispatch(logData gohf_extended.LogData) {
	logger.workerPool.Enqueue(logData)
}

func (logger *RabbitMQLogger) Close() {
	logger.workerPool.Close()
	logger.conn.Close()
}

func (logger *RabbitMQLogger) getChannel(goRoutineId int) (*amqp.Channel, error) {
	if logger.channels[goRoutineId] != nil {
		return logger.channels[goRoutineId], nil
	}

	ch, err := logger.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("conn.Channel error: %w", err)
	}

	logger.channels[goRoutineId] = ch
	return ch, nil
}
