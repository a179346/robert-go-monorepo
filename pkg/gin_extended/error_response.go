package gin_extended

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ztrue/tracerr"
)

var responseErrorDetail = false

func SetReponseErrorDetail(v bool) {
	responseErrorDetail = v
}

type ErrorResponseData struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

func (err ErrorResponseData) Error() string {
	return err.Message
}

type ErrorResponse struct {
	Status     int
	Message    string
	Err        error
	Unexpected bool

	bodyBytes []byte
}

func Error(c *gin.Context, statusCode int, message string, err error, unexpected bool) {
	withResponse(c, newErrorResponse(statusCode, message, err, unexpected))
}

func newErrorResponse(statusCode int, message string, err error, unexpected bool) *ErrorResponse {
	return &ErrorResponse{
		Status:     statusCode,
		Message:    message,
		Err:        err,
		Unexpected: unexpected,

		bodyBytes: nil,
	}
}

func (res ErrorResponse) Error() string {
	return fmt.Sprintf("http error %d: %s", res.Status, res.Message)
}

func (res ErrorResponse) Send(c *gin.Context) {
	if errors.Is(c.Request.Context().Err(), context.Canceled) {
		return
	}

	res.setHeader(c)
	c.Writer.WriteHeader(res.Status)
	//nolint:errcheck
	c.Writer.Write(res.getBodyBytes())
}

func (res *ErrorResponse) PrepareApiLog(c *gin.Context) (status int, bodyBytes []byte, logErr error, unexpected bool) {
	res.setHeader(c)
	return res.Status, res.getBodyBytes(), res.Err, res.Unexpected
}

func (res ErrorResponse) setHeader(c *gin.Context) {
	c.Header("Content-Type", "application/json")
}

func (res *ErrorResponse) getBodyBytes() []byte {
	if res.bodyBytes != nil {
		return res.bodyBytes
	}

	body := ErrorResponseData{
		Status:  res.Status,
		Message: res.Message,
	}
	if responseErrorDetail {
		body.Detail = tracerr.Sprint(res.Err)
	}
	res.bodyBytes, _ = json.Marshal(body)

	return res.bodyBytes
}
