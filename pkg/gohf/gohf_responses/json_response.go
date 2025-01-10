package gohf_responses

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/a179346/robert-go-monorepo/pkg/gohf"
)

type JsonResponse[T interface{}] struct {
	Status int
	Data   T
}

func NewJsonResponse[T interface{}](statusCode int, data T) JsonResponse[T] {
	return JsonResponse[T]{
		Status: statusCode,
		Data:   data,
	}
}

func (response JsonResponse[T]) Send(res gohf.ResponseWriter, req *gohf.Request) {
	if errors.Is(req.RootContext().Err(), context.Canceled) {
		return
	}

	res.SetHeader("Content-Type", "application/json")
	res.SetStatus(response.Status)
	json.NewEncoder(res).Encode(response.Data)
}