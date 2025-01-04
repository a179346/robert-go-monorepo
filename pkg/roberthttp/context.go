package roberthttp

type NextFunc func()

type Context struct {
	Res  Response
	Req  *Request
	Next NextFunc
}

func newContext(res Response, req *Request) *Context {
	return &Context{
		Res:  res,
		Req:  req,
		Next: func() {},
	}
}
