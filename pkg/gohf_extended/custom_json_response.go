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

	bodyBytes []byte
}

func NewCustomJsonResponse[T interface{}](statusCode int, data T) *CutsomJsonResponse[T] {
	bodyBytes, _ := json.Marshal(CutsomJsonResponseData[T]{data})
	return &CutsomJsonResponse[T]{
		Status: statusCode,

		bodyBytes: bodyBytes,
	}
}

func (res CutsomJsonResponse[T]) Send(w http.ResponseWriter, req *gohf.Request) {
	if errors.Is(req.RootContext().Err(), context.Canceled) {
		return
	}

	res.setHeader(w.Header())
	w.WriteHeader(res.Status)
	//nolint:errcheck
	w.Write(res.bodyBytes)
}

func (res *CutsomJsonResponse[T]) PrepareApiLog(header http.Header) (status int, bodyBytes []byte, logErr error, unexpected bool) {
	res.setHeader(header)
	return res.Status, res.bodyBytes, nil, false
}

func (res CutsomJsonResponse[T]) setHeader(header http.Header) {
	if header.Get("Content-Type") == "" {
		header.Set("Content-Type", "application/json")
	}
}
