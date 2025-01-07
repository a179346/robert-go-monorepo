package roberthttp_response

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/a179346/robert-go-monorepo/pkg/roberthttp"
)

type JsonHttpResponse[T interface{}] struct {
	Status int
	Data   T
}

func NewJsonResponse[T interface{}](statusCode int, data T) JsonHttpResponse[T] {
	return JsonHttpResponse[T]{
		Status: statusCode,
		Data:   data,
	}
}

func (r JsonHttpResponse[T]) Send(res roberthttp.Response, req *roberthttp.Request) {
	if errors.Is(req.RootContext().Err(), context.Canceled) {
		return
	}

	res.SetHeader("Content-Type", "application/json")
	res.SetStatus(r.Status)
	json.NewEncoder(res).Encode(r.Data)
}
