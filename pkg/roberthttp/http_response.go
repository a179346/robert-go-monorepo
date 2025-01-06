package roberthttp

type HttpResponse interface {
	Send(Response, *Request)
}
