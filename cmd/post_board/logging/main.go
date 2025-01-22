package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	post_board_config "github.com/a179346/robert-go-monorepo/internal/post_board/config"
	"github.com/a179346/robert-go-monorepo/pkg/console"
	"github.com/a179346/robert-go-monorepo/pkg/rabbitmq_consumerpool"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		console.Errorf("%s", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	go func() {
		<-ctx.Done()
		console.Info("Context cancelled. Logging system will gracefully shutdown.")
	}()

	rabbitMQConfig := post_board_config.GetRabbitMQConfig()

	conn, err := amqp.Dial(rabbitMQConfig.Url)
	if err != nil {
		return fmt.Errorf("amqp.Dial error: %w", err)
	}
	defer conn.Close()

	concurrency := post_board_config.GetLoggerConfig().ConsumerConcurrency
	consumerPool := rabbitmq_consumerpool.New(
		conn,
		&handlerImpl{},
		concurrency,
	)

	console.Infof("Logging system is serving. concurrency: %v", concurrency)
	consumerPool.Serve(ctx)
	return nil
}

type handlerImpl struct{}

func (handler *handlerImpl) Consume(ch *amqp.Channel) (<-chan amqp.Delivery, error) {
	return ch.Consume(
		post_board_config.GetLoggerConfig().ConsumerSourceQueue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
}

func (handler *handlerImpl) Handle(d amqp.Delivery) {
	console.Info(string(d.Body))
	//nolint:errcheck
	d.Ack(false)
}
