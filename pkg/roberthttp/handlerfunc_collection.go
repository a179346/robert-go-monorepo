package roberthttp

type HandlerFunc func(c *Context) HttpResponse

type HandlerFuncCollection struct {
	handlerFuncs []HandlerFunc
	pattern      string
	all          bool
}

func newHandlerFuncCollection(pattern string, all bool) HandlerFuncCollection {
	return HandlerFuncCollection{
		pattern: pattern,
		all:     all,
	}
}

func (h *HandlerFuncCollection) AddHandlerFuncs(handlerFuncs []HandlerFunc) {
	h.handlerFuncs = append(h.handlerFuncs, handlerFuncs...)
}

func (h *HandlerFuncCollection) GetHandlerFunc(idx int) HandlerFunc {
	return h.handlerFuncs[idx]
}

func (h *HandlerFuncCollection) Len() int {
	return len(h.handlerFuncs)
}
