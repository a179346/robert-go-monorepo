package rabbitmq_consumerpool

import (
	"context"
	"sync"
	"time"

	"github.com/a179346/robert-go-monorepo/pkg/console"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ConsumerPool struct {
	dial           DialFunc
	handlerFactory HandlerFactory
	concurrency    int
}

func New(dial DialFunc, handlerFactory HandlerFactory, concurrency int) *ConsumerPool {
	return &ConsumerPool{
		dial:           dial,
		handlerFactory: handlerFactory,
		concurrency:    concurrency,
	}
}

func (consumerPool *ConsumerPool) Serve(ctx context.Context) {
	for {
		if ctx.Err() != nil {
			return
		}
		consumerPool.dialAndServe(ctx)
		time.Sleep(1 * time.Second)
	}
}

func (consumerPool *ConsumerPool) dialAndServe(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	conn, err := consumerPool.dial()
	if err != nil {
		return
	}
	defer conn.Close()

	go func() {
		<-conn.NotifyClose(make(chan *amqp.Error))
		cancel()
	}()

	wg := new(sync.WaitGroup)

	wg.Add(consumerPool.concurrency)
	for i := 0; i < consumerPool.concurrency; i++ {
		consumer := newConsumer(conn, consumerPool.handlerFactory())
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
