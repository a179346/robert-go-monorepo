package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	post_board_config "github.com/a179346/robert-go-monorepo/internal/post_board/config"
	"github.com/a179346/robert-go-monorepo/pkg/console"
	"github.com/a179346/robert-go-monorepo/pkg/gohf_extended"
	"github.com/a179346/robert-go-monorepo/pkg/rabbitmq_consumerpool"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
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
	loggerConfig := post_board_config.GetLoggerConfig()

	cfg := elasticsearch.Config{
		Addresses: []string{loggerConfig.ElasticSearchAddress},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("elasticsearch.NewClient error: %w", err)
	}

	conn, err := amqp.Dial(rabbitMQConfig.Url)
	if err != nil {
		return fmt.Errorf("amqp.Dial error: %w", err)
	}
	defer conn.Close()

	concurrency := loggerConfig.ConsumerConcurrency
	consumerPool := rabbitmq_consumerpool.New(
		conn,
		&handlerImpl{
			sourceQueue: loggerConfig.ConsumerSourceQueue,
			es:          es,
		},
		concurrency,
	)

	console.Infof("Logging system is serving. concurrency: %v", concurrency)
	consumerPool.Serve(ctx)
	return nil
}

type handlerImpl struct {
	sourceQueue string
	es          *elasticsearch.Client
}

func (handler *handlerImpl) Consume(ch *amqp.Channel) (<-chan amqp.Delivery, error) {
	return ch.Consume(
		handler.sourceQueue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
}

func (handler *handlerImpl) Handle(d amqp.Delivery) {
	bodyBytes := d.Body

	var data gohf_extended.LogData
	err := json.Unmarshal(bodyBytes, &data)
	if err != nil {
		//nolint:errcheck
		d.Nack(false, false)
		return
	}

	startTime := time.UnixMilli(data.StartUnixMs)
	req := esapi.IndexRequest{
		Index:      fmt.Sprintf("post-board-api-%s", startTime.Format("20060102")),
		DocumentID: data.ID,
		Body:       bytes.NewReader(bodyBytes),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := req.Do(ctx, handler.es)
	if err != nil {
		//nolint:errcheck
		d.Nack(false, true)
		return
	}
	defer res.Body.Close()

	//nolint:errcheck
	d.Ack(false)
}
