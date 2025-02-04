package rabbitmqlogger

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/a179346/robert-go-monorepo/pkg/apilog"
	"github.com/a179346/robert-go-monorepo/pkg/workerpool"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQLogger struct {
	workerPool        *workerpool.WorkerPool[apilog.Data]
	getAmqpConnection func() (*amqp.Connection, error)
	channels          []*amqp.Channel
}

func New(
	getAmqpConnection func() (*amqp.Connection, error),
	exchange string,
) *RabbitMQLogger {
	const concurrency = 8
	retryDelays := []time.Duration{2, 2, 2, 4, 6, 8, 10, 12, 14, 20}

	logger := &RabbitMQLogger{
		getAmqpConnection: getAmqpConnection,
		channels:          make([]*amqp.Channel, concurrency),
	}

	workerPool := workerpool.New(func(data apilog.Data, goRoutineId int) {
		body, err := json.Marshal(data)
		if err != nil {
			return
		}

		for _, retryDelay := range retryDelays {
			ch, err := logger.getChannel(goRoutineId)
			if err == nil {
				err = ch.Publish(
					exchange,
					"",
					false,
					false,
					amqp.Publishing{
						MessageId:    data.ID,
						AppId:        data.App,
						ContentType:  "application/json",
						DeliveryMode: 2,
						Body:         body,
						Timestamp:    time.UnixMilli(data.StartUnixMs),
					},
				)
			}

			if err == nil {
				return
			}
			time.Sleep(retryDelay * time.Second)
		}
	}, concurrency, 1024, true)
	logger.workerPool = workerPool

	return logger
}

func (logger *RabbitMQLogger) Dispatch(data apilog.Data) {
	logger.workerPool.Enqueue(data)
}

func (logger *RabbitMQLogger) Close() {
	logger.workerPool.Close()
	if conn, err := logger.getAmqpConnection(); err == nil {
		conn.Close()
	}
}

func (logger *RabbitMQLogger) getChannel(goRoutineId int) (*amqp.Channel, error) {
	if ch := logger.channels[goRoutineId]; ch != nil && !ch.IsClosed() {
		return ch, nil
	}

	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	if ch := logger.channels[goRoutineId]; ch != nil && !ch.IsClosed() {
		return ch, nil
	}

	conn, err := logger.getAmqpConnection()
	if err != nil {
		return nil, fmt.Errorf("logger.getAmqpConnection error: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("conn.Channel error: %w", err)
	}

	logger.channels[goRoutineId] = ch
	return ch, nil
}
