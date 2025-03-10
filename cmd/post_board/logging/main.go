package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/a179346/robert-go-monorepo/pkg/console"
	"github.com/a179346/robert-go-monorepo/pkg/es_bulkrequester"
	"github.com/a179346/robert-go-monorepo/pkg/rabbitmq_consumerpool"
	post_board_config "github.com/a179346/robert-go-monorepo/services/post_board/config"
	"github.com/a179346/robert-go-monorepo/services/post_board/logging_helper"
	"github.com/elastic/go-elasticsearch/v8"
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
	loggingConfig := post_board_config.GetLoggingConfig()

	cfg := elasticsearch.Config{
		Addresses: []string{loggingConfig.ElasticSearchAddress},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("elasticsearch.NewClient error: %w", err)
	}

	concurrency := loggingConfig.ConsumerConcurrency
	consumerPool := rabbitmq_consumerpool.New(
		func() (*amqp.Connection, error) {
			return amqp.Dial(rabbitMQConfig.Url)
		},
		func() rabbitmq_consumerpool.Handler {
			return logging_helper.NewHandler(
				loggingConfig.ConsumerSourceQueue,
				loggingConfig.ElasticSearchIndexPrefix,
				es_bulkrequester.New(es, 500, 30*time.Second),
			)
		},
		concurrency,
	)

	console.Infof("Logging system is serving. concurrency: %v", concurrency)
	consumerPool.Serve(ctx)
	return nil
}
