package gohf

import (
	"net/http"
)

type prefixedHandlerFunc struct {
	prefix string
	f      HandlerFunc
}

func newPrefixedHandlerFunc(prefix string, f HandlerFunc) *prefixedHandlerFunc {
	return &prefixedHandlerFunc{
		prefix: prefix,
		f:      f,
	}
}

type httpHandler struct {
	fullPrefix           string
	prefixedHandlerFuncs []*prefixedHandlerFunc
}

func newHttpHandler(fullPrefix string) *httpHandler {
	return &httpHandler{
		fullPrefix: fullPrefix,
	}
}

func (httpHandler *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := newResponseWriter(w)
	req := newRequest(r, httpHandler.fullPrefix)
	c := newContext(res, req)

	var handle func(idx int) Response
	handle = func(idx int) Response {
		if idx == len(httpHandler.prefixedHandlerFuncs) {
			return nil
		}

		c.Next = func() Response { return handle(idx + 1) }

		handler := httpHandler.prefixedHandlerFuncs[idx]

		originalPrefix := req.getCurrPrefix()
		req.setCurrPrefix(handler.prefix)
		response := handler.f(c)
		req.setCurrPrefix(originalPrefix)

		return response
	}

	if response := handle(0); response != nil {
		req.setCurrPrefix("")
		response.Send(c.Res, c.Req)
	}
}

func (httpHandler *httpHandler) addHandlerFunc(prefix string, f HandlerFunc) {
	httpHandler.prefixedHandlerFuncs = append(
		httpHandler.prefixedHandlerFuncs,
		newPrefixedHandlerFunc(prefix, f),
	)
}

func (httpHandler *httpHandler) len() int {
	return len(httpHandler.prefixedHandlerFuncs)
}
