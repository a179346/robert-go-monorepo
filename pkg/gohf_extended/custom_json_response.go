package gohf_extended

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gohf-http/gohf/v6"
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

func (res CutsomJsonResponse[T]) Send(w http.ResponseWriter, req *gohf.Request) {
	if errors.Is(req.RootContext().Err(), context.Canceled) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.Status)
	//nolint:errcheck
	json.NewEncoder(w).Encode(res.Data)
}
