package rabbitmq_consumerpool

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn    *amqp.Connection
	handler Handler
}

func newConsumer(conn *amqp.Connection, handler Handler) *Consumer {
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
		case <-ctx.Done():
			return nil

		case d := <-msgs:
			consumer.handler.Handle(d)
		}
	}
}
