package rabbitmq_consumerpool

import (
	"context"
	"sync"

	"github.com/a179346/robert-go-monorepo/pkg/console"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ConsumerPool struct {
	conn        *amqp.Connection
	handler     Handler
	concurrency int
}

func New(conn *amqp.Connection, handler Handler, concurrency int) *ConsumerPool {
	return &ConsumerPool{
		conn:        conn,
		handler:     handler,
		concurrency: concurrency,
	}
}

func (consumerPool *ConsumerPool) Serve(ctx context.Context) {
	wg := new(sync.WaitGroup)

	wg.Add(consumerPool.concurrency)
	for i := 0; i < consumerPool.concurrency; i++ {
		consumer := newConsumer(consumerPool.conn, consumerPool.handler)
		go func(consumer *Consumer) {
			err := consumer.Serve(ctx)
			if err != nil {
				console.Errorf("consumer.Serve error: %v", err)
			}
			wg.Done()
		}(consumer)
	}

	wg.Wait()
}
