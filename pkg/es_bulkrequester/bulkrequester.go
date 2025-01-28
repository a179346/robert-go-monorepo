package es_bulkrequester

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

type BulkRequester struct {
	es             *elasticsearch.Client
	maxBatchSize   int
	mu             *sync.Mutex
	batchFull      chan<- struct{}
	buf            *bytes.Buffer
	bulkItemEvents []*BulkItemEvent
	closing        chan<- struct{}
	closed         <-chan struct{}
}

func New(es *elasticsearch.Client, maxBatchSize int, writeInterval time.Duration) *BulkRequester {
	batchFull := make(chan struct{})
	closing := make(chan struct{})
	closed := make(chan struct{})

	requester := &BulkRequester{
		es:             es,
		maxBatchSize:   maxBatchSize,
		mu:             new(sync.Mutex),
		batchFull:      batchFull,
		buf:            new(bytes.Buffer),
		bulkItemEvents: nil,
		closing:        closing,
		closed:         closed,
	}

	go func() {
		defer close(closed)

		for {
			select {
			case <-closing:
				requester.sendRequest()
				return
			default:
			}

			select {
			case <-time.After(writeInterval):
				requester.sendRequest()
			case <-batchFull:
				requester.sendRequest()
			case <-closing:
				requester.sendRequest()
				return
			}
		}
	}()

	return requester
}

func (requester *BulkRequester) Close() {
	requester.closing <- struct{}{}
	<-requester.closed
}

func (requester *BulkRequester) AddRequest(meta []byte, data []byte, bulkItemEvent *BulkItemEvent) {
	requester.mu.Lock()

	requester.buf.Grow(len(meta) + len(data))
	requester.buf.Write(meta)
	requester.buf.Write(data)
	requester.bulkItemEvents = append(requester.bulkItemEvents, bulkItemEvent)

	requester.mu.Unlock()

	if len(requester.bulkItemEvents) >= requester.maxBatchSize {
		requester.batchFull <- struct{}{}
	}
}

func (requester *BulkRequester) sendRequest() {
	requester.mu.Lock()
	buf := requester.buf
	requester.buf = new(bytes.Buffer)
	bulkItemEvents := requester.bulkItemEvents
	requester.bulkItemEvents = nil
	requester.mu.Unlock()

	if len(bulkItemEvents) == 0 {
		return
	}

	bulkReq := esapi.BulkRequest{
		Body: buf,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := bulkReq.Do(ctx, requester.es)
	if err != nil || res.IsError() {
		for _, bulkItemEvent := range bulkItemEvents {
			bulkItemEvent.onError()
		}
		return
	}
	defer res.Body.Close()

	var blk *bulkResponse
	if err := json.NewDecoder(res.Body).Decode(&blk); err != nil {
		for _, bulkItemEvent := range bulkItemEvents {
			bulkItemEvent.onError()
		}
		return
	}

	for i, d := range blk.Items {
		if d.Index.Status > 201 {
			bulkItemEvents[i].onError()
		} else {
			bulkItemEvents[i].onSuccess()
		}
	}
}
