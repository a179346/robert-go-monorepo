package roberthttp

import (
	"context"
	"mime/multipart"
	"net/http"
	"net/url"
)

type Request struct {
	req *http.Request
}

func newRequest(req *http.Request) *Request {
	return &Request{req}
}

func (req *Request) PathValue(name string) string {
	return req.req.PathValue(name)
}

func (req *Request) URL() *url.URL {
	return req.req.URL
}

func (req *Request) Context() context.Context {
	return req.req.Context()
}

func (req *Request) FormFile(key string) (multipart.File, *multipart.FileHeader, error) {
	return req.req.FormFile(key)
}

func (req *Request) FormValue(key string) string {
	return req.req.FormValue(key)
}

func (req *Request) GetHeader(key string) string {
	return req.req.Header.Get(key)
}