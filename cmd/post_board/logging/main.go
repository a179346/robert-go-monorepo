package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
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
	loggingConfig := post_board_config.GetLoggingConfig()

	cfg := elasticsearch.Config{
		Addresses: []string{loggingConfig.ElasticSearchAddress},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("elasticsearch.NewClient error: %w", err)
	}

	_ = upsertILMLifecycle(ctx, es)
	_ = upsertIndexTemplate(ctx, es)

	concurrency := loggingConfig.ConsumerConcurrency
	consumerPool := rabbitmq_consumerpool.New(
		&handlerImpl{
			url:         rabbitMQConfig.Url,
			sourceQueue: loggingConfig.ConsumerSourceQueue,
			es:          es,
			indexPrefix: loggingConfig.ElasticSearchIndexPrefix,
		},
		concurrency,
	)

	console.Infof("Logging system is serving. concurrency: %v", concurrency)
	consumerPool.Serve(ctx)
	return nil
}

func upsertILMLifecycle(ctx context.Context, es *elasticsearch.Client) error {
	prefix := post_board_config.GetLoggingConfig().ElasticSearchIndexPrefix
	policy := prefix + "30-days"
	body := `{"policy":{"phases":{"hot":{"min_age":"0ms","actions":{}},"warm":{"min_age":"2d","actions":{}},"delete":{"min_age":"30d","actions":{"delete":{"delete_searchable_snapshot":true}}}},"deprecated":false}}`
	upsertILMLifecycleReq := esapi.ILMPutLifecycleRequest{
		Policy: policy,
		Body:   strings.NewReader(body),
	}
	res, err := upsertILMLifecycleReq.Do(ctx, es)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

func upsertIndexTemplate(ctx context.Context, es *elasticsearch.Client) error {
	prefix := post_board_config.GetLoggingConfig().ElasticSearchIndexPrefix
	policy := prefix + "30-days"
	indexPattern := prefix + "*"
	body := fmt.Sprintf(
		`{"template":{"settings":{"index":{"lifecycle":{"name":"%s"}}}},"index_patterns":["%s"],"data_stream":{}}`,
		policy,
		indexPattern,
	)
	upsertIndexTemplateReq := esapi.IndicesPutIndexTemplateRequest{
		Name: post_board_config.GetLoggingConfig().ElasticSearchIndexPrefix,
		Body: strings.NewReader(body),
	}
	res, err := upsertIndexTemplateReq.Do(ctx, es)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

type handlerImpl struct {
	url         string
	sourceQueue string
	es          *elasticsearch.Client
	indexPrefix string
}

func (handler *handlerImpl) Dial() (*amqp.Connection, error) {
	return amqp.Dial(handler.url)
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

	var data gohf_extended.ApiLogData
	err := json.Unmarshal(bodyBytes, &data)
	if err != nil {
		//nolint:errcheck
		d.Nack(false, false)
		return
	}

	timestamp := time.UnixMilli(data.Timestamp)
	req := esapi.IndexRequest{
		Index:      handler.indexPrefix + timestamp.Format("20060102"),
		DocumentID: data.ID,
		OpType:     "create",
		Body:       bytes.NewReader(bodyBytes),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := req.Do(ctx, handler.es)
	if err != nil || res.StatusCode >= 400 {
		//nolint:errcheck
		d.Nack(false, true)
		return
	}
	defer res.Body.Close()

	//nolint:errcheck
	d.Ack(false)
}
