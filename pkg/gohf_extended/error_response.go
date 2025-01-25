package gohf_extended

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gohf-http/gohf/v6"
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

func NewErrorResponse(statusCode int, message string, err error, unexpected bool) *ErrorResponse {
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

func (res ErrorResponse) Send(w http.ResponseWriter, req *gohf.Request) {
	if errors.Is(req.RootContext().Err(), context.Canceled) {
		return
	}

	res.setHeader(w.Header())
	w.WriteHeader(res.Status)
	//nolint:errcheck
	w.Write(res.getBodyBytes())
}

func (res *ErrorResponse) PrepareApiLog(header http.Header) (status int, bodyBytes []byte, logErr error, unexpected bool) {
	res.setHeader(header)
	return res.Status, res.getBodyBytes(), res.Err, res.Unexpected
}

func (res ErrorResponse) setHeader(header http.Header) {
	if header.Get("Content-Type") == "" {
		header.Set("Content-Type", "application/json")
	}
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
