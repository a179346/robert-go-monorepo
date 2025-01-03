package roberthttp

import (
	"encoding/json"
	"io"
	"net/http"
)

type Response struct {
	w http.ResponseWriter
}

type ErrorResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func newResponse(w http.ResponseWriter) Response {
	return Response{w}
}

func (res Response) SetHeader(key, value string) {
	res.w.Header().Set(key, value)
}

func (res Response) SetStatue(statusCode int) {
	res.w.WriteHeader(statusCode)
}

func (res Response) WriteJson(statusCode int, data interface{}) error {
	res.SetHeader("Content-Type", "application/json")
	res.w.WriteHeader(statusCode)
	return json.NewEncoder(res.w).Encode(data)
}

func (res Response) WriteError(statusCode int, message string, data interface{}) error {
	return res.WriteJson(statusCode, ErrorResponse{
		Message: message,
		Data:    data,
	})
}

func (res Response) GetWriter() io.Writer {
	return res.w
}
