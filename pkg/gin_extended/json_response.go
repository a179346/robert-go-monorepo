package gin_extended

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/gin-gonic/gin"
)

type JsonResponseData[T interface{}] struct {
	Data T `json:"data"`
}

type JsonResponse[T interface{}] struct {
	Status int
	Data   JsonResponseData[T]

	bodyBytes []byte
}

func JSON[T interface{}](c *gin.Context, statusCode int, data T) {
	withResponse(c, newJsonResponse(statusCode, data))
}

func newJsonResponse[T interface{}](statusCode int, data T) *JsonResponse[T] {
	return &JsonResponse[T]{
		Status: statusCode,
		Data:   JsonResponseData[T]{data},

		bodyBytes: nil,
	}
}

func (res JsonResponse[T]) Send(c *gin.Context) {
	if errors.Is(c.Request.Context().Err(), context.Canceled) {
		return
	}

	res.setHeader(c)
	c.Writer.WriteHeader(res.Status)
	//nolint:errcheck
	c.Writer.Write(res.getBodyBytes())
}

func (res *JsonResponse[T]) PrepareApiLog(c *gin.Context) (status int, bodyBytes []byte, logErr error, unexpected bool) {
	res.setHeader(c)
	return res.Status, res.getBodyBytes(), nil, false
}

func (res JsonResponse[T]) setHeader(c *gin.Context) {
	c.Header("Content-Type", "application/json")
}

func (res *JsonResponse[T]) getBodyBytes() []byte {
	if res.bodyBytes != nil {
		return res.bodyBytes
	}

	res.bodyBytes, _ = json.Marshal(res.Data)
	return res.bodyBytes
}
