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

	bodyBytes []byte
}

func JSON[T interface{}](c *gin.Context, statusCode int, data T) {
	withResponse(c, newJsonResponse(statusCode, data))
}

func newJsonResponse[T interface{}](statusCode int, data T) *JsonResponse[T] {
	bodyBytes, _ := json.Marshal(JsonResponseData[T]{data})

	return &JsonResponse[T]{
		Status: statusCode,

		bodyBytes: bodyBytes,
	}
}

func (res JsonResponse[T]) Send(c *gin.Context) {
	if errors.Is(c.Request.Context().Err(), context.Canceled) {
		return
	}

	res.setHeader(c)
	c.Writer.WriteHeader(res.Status)
	//nolint:errcheck
	c.Writer.Write(res.bodyBytes)
}

func (res *JsonResponse[T]) PrepareApiLog(c *gin.Context) (status int, bodyBytes []byte, logErr error, unexpected bool) {
	res.setHeader(c)
	return res.Status, res.bodyBytes, nil, false
}

func (res JsonResponse[T]) setHeader(c *gin.Context) {
	c.Header("Content-Type", "application/json")
}
