package roberthttp

import (
	"io"
	"net/http"
	"time"
)

type Response struct {
	w http.ResponseWriter
}

func newResponse(w http.ResponseWriter) Response {
	return Response{
		w: w,
	}
}

func (res Response) SetHeader(key, value string) {
	res.w.Header().Set(key, value)
}

func (res Response) SetStatus(statusCode int) {
	res.w.WriteHeader(statusCode)
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
