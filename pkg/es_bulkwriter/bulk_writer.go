package es_bulkwriter

import (
	"bytes"
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type bulkResponse struct {
	Errors bool `json:"errors"`
	Items  []struct {
		Index struct {
			ID     string `json:"_id"`
			Result string `json:"result"`
			Status int    `json:"status"`
			Error  struct {
				Type   string `json:"type"`
				Reason string `json:"reason"`
				Cause  struct {
					Type   string `json:"type"`
					Reason string `json:"reason"`
				} `json:"caused_by"`
			} `json:"error"`
		} `json:"index"`
	} `json:"items"`
}

type BulkWriter struct {
	es              *elasticsearch.Client
	maxBatchSize    int
	mu              *sync.Mutex
	batchFull       chan<- struct{}
	buf             *bytes.Buffer
	bulkWriteEvents []*BulkWriteEvent
	closing         chan<- struct{}
	closed          chan struct{}
}

func New(es *elasticsearch.Client, maxBatchSize int, writeInterval time.Duration) *BulkWriter {
	batchFull := make(chan struct{})
	closing := make(chan struct{})
	closed := make(chan struct{})

	writer := &BulkWriter{
		es:              es,
		maxBatchSize:    maxBatchSize,
		mu:              new(sync.Mutex),
		batchFull:       batchFull,
		buf:             new(bytes.Buffer),
		bulkWriteEvents: nil,
		closing:         closing,
		closed:          closed,
	}

	go func() {
		defer close(closed)

		for {
			select {
			case <-time.After(writeInterval):
				writer.sendRequest()
			case <-batchFull:
				writer.sendRequest()
			case <-closing:
				writer.sendRequest()
				return
			}
		}
	}()

	return writer
}

func (bulkWriter *BulkWriter) Close() {
	bulkWriter.closing <- struct{}{}
	<-bulkWriter.closed
}

func (bulkWriter *BulkWriter) AddRequest(meta []byte, data []byte, bulkWriteEvent *BulkWriteEvent) {
	bulkWriter.mu.Lock()

	bulkWriter.buf.Grow(len(meta) + len(data))
	bulkWriter.buf.Write(meta)
	bulkWriter.buf.Write(data)
	bulkWriter.bulkWriteEvents = append(bulkWriter.bulkWriteEvents, bulkWriteEvent)

	bulkWriter.mu.Unlock()

	if len(bulkWriter.bulkWriteEvents) >= bulkWriter.maxBatchSize {
		bulkWriter.batchFull <- struct{}{}
	}
}

func (bulkWriter *BulkWriter) sendRequest() {
	bulkWriter.mu.Lock()
	buf := bulkWriter.buf
	bulkWriter.buf = new(bytes.Buffer)
	bulkWriteEvents := bulkWriter.bulkWriteEvents
	bulkWriter.bulkWriteEvents = nil
	bulkWriter.mu.Unlock()

	if len(bulkWriteEvents) == 0 {
		return
	}

	bulkReq := esapi.BulkRequest{
		Body: buf,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := bulkReq.Do(ctx, bulkWriter.es)
	if err != nil || res.IsError() {
		for _, bulkWriteEvent := range bulkWriteEvents {
			bulkWriteEvent.onError()
		}
		return
	}
	defer res.Body.Close()

	var blk *bulkResponse
	if err := json.NewDecoder(res.Body).Decode(&blk); err != nil {
		for _, bulkWriteEvent := range bulkWriteEvents {
			bulkWriteEvent.onError()
		}
		return
	}

	for i, d := range blk.Items {
		if d.Index.Status > 201 {
			bulkWriteEvents[i].onError()
		} else {
			bulkWriteEvents[i].onSuccess()
		}
	}
}
