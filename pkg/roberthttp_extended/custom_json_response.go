package roberthttp_extended

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/a179346/robert-go-monorepo/pkg/roberthttp"
)

type CutsomJsonResponseData[T interface{}] struct {
	Data T `json:"data"`
}

type CutsomJsonHttpResponse[T interface{}] struct {
	Status int
	Data   CutsomJsonResponseData[T]
}

func NewCustomJsonResponse[T interface{}](statusCode int, data T) CutsomJsonHttpResponse[T] {
	return CutsomJsonHttpResponse[T]{
		Status: statusCode,
		Data:   CutsomJsonResponseData[T]{data},
	}
}

func (r CutsomJsonHttpResponse[T]) Send(res roberthttp.Response, req *roberthttp.Request) {
	if errors.Is(req.RootContext().Err(), context.Canceled) {
		return
	}

	res.SetHeader("Content-Type", "application/json")
	res.SetStatus(r.Status)
	json.NewEncoder(res).Encode(r.Data)
}
