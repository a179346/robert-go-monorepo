package roberthttp

type HttpResponse interface {
	Send(Response, *Request)
}

type HandlerFunc func(c *Context) HttpResponse

type NextFunc func() HttpResponse

type Context struct {
	Res  Response
	Req  *Request
	Next NextFunc
}

func newContext(res Response, req *Request) *Context {
	return &Context{
		Res:  res,
		Req:  req,
		Next: func() HttpResponse { return nil },
	}
}
