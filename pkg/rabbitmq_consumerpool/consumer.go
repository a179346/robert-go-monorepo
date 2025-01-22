package rabbitmq_consumerpool

import (
	"context"
	"time"

	ampq "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn    *ampq.Connection
	handler Handler
}

func newConsumer(conn *ampq.Connection, handler Handler) *Consumer {
	return &Consumer{
		conn:    conn,
		handler: handler,
	}
}

func (consumer *Consumer) Serve(ctx context.Context) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	msgs, err := consumer.handler.Consume(ch)
	if err != nil {
		return err
	}

	for {
		if ctx.Err() != nil {
			return nil
		}

		select {
		case d := <-msgs:
			consumer.handler.Handle(d)

		case <-time.After(20 * time.Millisecond):
		}
	}
}
