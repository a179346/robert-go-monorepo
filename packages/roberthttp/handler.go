package roberthttp

import (
	"net/http"
)

type HandlerFunc func(c *Context)

type Handler interface {
	AddHandlerFuncs([]HandlerFunc)
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type HandlerImpl struct {
	handlerFuncs []HandlerFunc
}

func newHandler() HandlerImpl {
	return HandlerImpl{}
}

func (h *HandlerImpl) AddHandlerFuncs(handlerFuncs []HandlerFunc) {
	h.handlerFuncs = append(h.handlerFuncs, handlerFuncs...)
}

func (handler HandlerImpl) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	var handle func(idx int)
	handle = func(idx int) {
		if idx == len(handler.handlerFuncs) {
			return
		}
		c.Next = func() { handle(idx + 1) }

		handler.handlerFuncs[idx](c)
	}

	handle(0)
}

type PatternHandler struct {
	HandlerImpl
	pattern string
}

func newPatternHandler(pattern string) PatternHandler {
	return PatternHandler{
		HandlerImpl: newHandler(),
		pattern:     pattern,
	}
}

type AllHandler struct {
	HandlerImpl
}

func newAllHandler() AllHandler {
	return AllHandler{newHandler()}
}
