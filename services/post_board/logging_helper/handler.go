package logging_helper

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/a179346/robert-go-monorepo/pkg/apilog"
	"github.com/a179346/robert-go-monorepo/pkg/es_bulkrequester"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Handler struct {
	sourceQueue   string
	indexPrefix   string
	bulkRequester *es_bulkrequester.BulkRequester
}

func NewHandler(
	sourceQueue string,
	indexPrefix string,
	bulkRequester *es_bulkrequester.BulkRequester,
) *Handler {
	return &Handler{
		sourceQueue:   sourceQueue,
		indexPrefix:   indexPrefix,
		bulkRequester: bulkRequester,
	}
}

func (handler *Handler) Consume(ch *amqp.Channel) (<-chan amqp.Delivery, error) {
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

func (handler *Handler) Handle(d amqp.Delivery) {
	bodyBytes := d.Body

	var data apilog.Data
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		//nolint:errcheck
		d.Nack(false, false)
		return
	}

	timestamp := time.UnixMilli(data.Timestamp)
	index := handler.indexPrefix + timestamp.Format("20060102")

	meta := []byte(fmt.Sprintf(`{"create":{"_index":"%s","_id":"%s"}}%s`, index, data.ID, "\n"))
	bodyBytes = append(bodyBytes, "\n"...)

	handler.bulkRequester.AddRequest(
		meta,
		bodyBytes,
		es_bulkrequester.NewBulkItemEvent(
			func() { _ = d.Ack(false) },
			func() { _ = d.Nack(false, true) },
		),
	)
}

func (handler *Handler) Close() {
	handler.bulkRequester.Close()
}
