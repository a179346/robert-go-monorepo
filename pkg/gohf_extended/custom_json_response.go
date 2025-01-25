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

	bodyBytes []byte
}

func NewCustomJsonResponse[T interface{}](statusCode int, data T) *CutsomJsonResponse[T] {
	return &CutsomJsonResponse[T]{
		Status: statusCode,
		Data:   CutsomJsonResponseData[T]{data},

		bodyBytes: nil,
	}
}

func (res CutsomJsonResponse[T]) Send(w http.ResponseWriter, req *gohf.Request) {
	if errors.Is(req.RootContext().Err(), context.Canceled) {
		return
	}

	res.setHeader(w.Header())
	w.WriteHeader(res.Status)
	//nolint:errcheck
	w.Write(res.getBodyBytes())
}

func (res *CutsomJsonResponse[T]) PrepareApiLog(header http.Header) (status int, bodyBytes []byte, logErr error, unexpected bool) {
	res.setHeader(header)
	return res.Status, res.getBodyBytes(), nil, false
}

func (res CutsomJsonResponse[T]) setHeader(header http.Header) {
	if header.Get("Content-Type") == "" {
		header.Set("Content-Type", "application/json")
	}
}

func (res *CutsomJsonResponse[T]) getBodyBytes() []byte {
	if res.bodyBytes != nil {
		return res.bodyBytes
	}

	res.bodyBytes, _ = json.Marshal(res.Data)
	return res.bodyBytes
}
