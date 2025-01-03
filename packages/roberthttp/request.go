package roberthttp

import (
	"context"
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
