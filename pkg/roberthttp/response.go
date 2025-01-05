package roberthttp

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type Response struct {
	w       http.ResponseWriter
	options *ResponseOptions
}

func newResponse(w http.ResponseWriter, options *ResponseOptions) Response {
	return Response{
		w:       w,
		options: options,
	}
}

func (res Response) SetHeader(key, value string) {
	res.w.Header().Set(key, value)
}

func (res Response) SetStatus(statusCode int) {
	res.w.WriteHeader(statusCode)
}

func (res Response) WriteJson(statusCode int, data interface{}) error {
	res.SetHeader("Content-Type", "application/json")
	res.w.WriteHeader(statusCode)

	var v interface{}
	if res.options.JsonWrapper != nil {
		v = res.options.JsonWrapper(statusCode, data)
	} else {
		v = defaultResponseJsonWrapper(statusCode, data)
	}

	return json.NewEncoder(res.w).Encode(v)
}

func (res Response) WriteError(statusCode int, message string, info interface{}) error {
	res.SetHeader("Content-Type", "application/json")
	res.w.WriteHeader(statusCode)

	var v interface{}
	if res.options.ErrorWrapper != nil {
		v = res.options.ErrorWrapper(statusCode, message, info)
	} else {
		v = defaultResponseErrorWrapper(statusCode, message, info)
	}

	return json.NewEncoder(res.w).Encode(v)
}

func (res Response) ServeFile(req *Request, filepath string) {
	http.ServeFile(res.w, req.req, filepath)
}

func (res Response) ServeContent(req *Request, name string, modtime time.Time, content io.ReadSeeker) {
	http.ServeContent(res.w, req.req, name, modtime, content)
}

func (res Response) GetWriter() io.Writer {
	return res.w
}
