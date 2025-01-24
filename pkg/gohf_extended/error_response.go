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
}

func NewErrorResponse(statusCode int, message string, err error, unexpected bool) ErrorResponse {
	return ErrorResponse{
		Status:     statusCode,
		Message:    message,
		Err:        err,
		Unexpected: unexpected,
	}
}

func (res ErrorResponse) Error() string {
	return fmt.Sprintf("http error %d: %s", res.Status, res.Message)
}

func (res ErrorResponse) Send(w http.ResponseWriter, req *gohf.Request) {
	w.Header().Set("Content-Type", "application/json")

	body := ErrorResponseData{
		Status:  res.Status,
		Message: res.Message,
	}
	if responseErrorDetail {
		body.Detail = tracerr.Sprint(res.Err)
	}
	bodyBytes, _ := json.Marshal(body)

	if apiLogger != nil {
		log(w, req, res.Status, bodyBytes, res.Err, res.Unexpected)
	}

	if errors.Is(req.RootContext().Err(), context.Canceled) {
		return
	}

	w.WriteHeader(res.Status)
	//nolint:errcheck
	w.Write(bodyBytes)
}
