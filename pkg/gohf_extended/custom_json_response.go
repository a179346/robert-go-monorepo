package gohf_extended

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/a179346/robert-go-monorepo/pkg/gohf"
)

type CutsomJsonResponseData[T interface{}] struct {
	Data T `json:"data"`
}

type CutsomJsonResponse[T interface{}] struct {
	Status int
	Data   CutsomJsonResponseData[T]
}

func NewCustomJsonResponse[T interface{}](statusCode int, data T) CutsomJsonResponse[T] {
	return CutsomJsonResponse[T]{
		Status: statusCode,
		Data:   CutsomJsonResponseData[T]{data},
	}
}

func (response CutsomJsonResponse[T]) Send(res gohf.ResponseWriter, req *gohf.Request) {
	if errors.Is(req.RootContext().Err(), context.Canceled) {
		return
	}

	res.SetHeader("Content-Type", "application/json")
	res.SetStatus(response.Status)
	json.NewEncoder(res).Encode(response.Data)
}