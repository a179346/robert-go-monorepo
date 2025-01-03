package roberthttp

import "net/http"

type NextFunc func()

type Context struct {
	Res  Response
	Req  *Request
	Next NextFunc
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Res:  newResponse(w),
		Req:  newRequest(req),
		Next: func() {},
	}
}
