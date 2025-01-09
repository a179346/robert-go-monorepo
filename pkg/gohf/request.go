package gohf

import (
	"context"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

type Request struct {
	req        *http.Request
	ctx        context.Context
	timestamp  time.Time
	fullPrefix string
	currPrefix string
}

func newRequest(req *http.Request, fullPrefix string) *Request {
	return &Request{
		req:        req,
		ctx:        req.Context(),
		timestamp:  time.Now(),
		fullPrefix: fullPrefix,
	}
}

func (req *Request) GetTimestamp() time.Time {
	return req.timestamp
}

func (req *Request) Host() string {
	return req.req.Host
}

func (req *Request) Path() string {
	return strings.TrimPrefix(req.FullPath(), req.currPrefix)
}

func (req *Request) FullPath() string {
	return req.fullPrefix + req.req.URL.Path
}

func (req *Request) RootContext() context.Context {
	return req.req.Context()
}

func (req *Request) Context() context.Context {
	return req.ctx
}

func (req *Request) SetContext(ctx context.Context) {
	req.ctx = ctx
}

func (req *Request) GetHeader(key string) string {
	return req.req.Header.Get(key)
}

func (req *Request) PathValue(name string) string {
	return req.req.PathValue(name)
}

func (req *Request) GetQuery(key string) string {
	return req.req.URL.Query().Get(key)
}

func (req *Request) FormFile(key string) (multipart.File, *multipart.FileHeader, error) {
	return req.req.FormFile(key)
}

func (req *Request) FormValue(key string) string {
	return req.req.FormValue(key)
}

func (req *Request) getCurrPrefix() string {
	return req.currPrefix
}

func (req *Request) setCurrPrefix(prefix string) {
	req.currPrefix = prefix
}
