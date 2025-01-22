package rabbitmq_consumerpool

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Handler interface {
	Consume(ch *amqp.Channel) (<-chan amqp.Delivery, error)
	Handle(d amqp.Delivery)
}
