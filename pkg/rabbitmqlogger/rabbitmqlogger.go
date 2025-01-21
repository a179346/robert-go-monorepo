package rabbitmqlogger

import (
	"context"
	"fmt"
	"time"

	"github.com/a179346/robert-go-monorepo/pkg/flushworker"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQLogger struct {
	worker   *flushworker.FlushWorker[[]byte]
	conn     *amqp.Connection
	channels []*amqp.Channel
}

func New(conn *amqp.Connection, exchange string) *RabbitMQLogger {
	const concurrency = 8
	const maxRetryCnt = 5

	logger := &RabbitMQLogger{
		conn:     conn,
		channels: make([]*amqp.Channel, concurrency),
	}

	worker := flushworker.New(func(v []byte, goRoutineId int) {
		for retryCnt := 1; retryCnt <= maxRetryCnt; retryCnt++ {
			ch, err := logger.getChannel(goRoutineId)
			if err == nil {
				err = ch.Publish(
					exchange,
					"",
					false,
					false,
					amqp.Publishing{
						DeliveryMode: 2,
						Body:         v,
					},
				)
			}

			if err == nil {
				return
			}
			time.Sleep(time.Duration(retryCnt*2) * time.Second)
		}
	}, concurrency, 1024)
	logger.worker = worker

	return logger
}

func (logger *RabbitMQLogger) Write(v []byte) {
	logger.worker.AddJob(v)
}

func (logger *RabbitMQLogger) Close(ctx context.Context) {
	logger.worker.Close(ctx)
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
