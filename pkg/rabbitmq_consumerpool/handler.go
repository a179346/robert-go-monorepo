package rabbitmq_consumerpool

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Handler interface {
	Dial() (*amqp.Connection, error)
	Consume(ch *amqp.Channel) (<-chan amqp.Delivery, error)
	Handle(d amqp.Delivery)
}
