package roberthttp

type HandlerFunc func(c *Context)

type HandlerFuncCollection interface {
	AddHandlerFuncs([]HandlerFunc)
	GetHandlerFunc(int) HandlerFunc
	Len() int
}

type HandlerFuncCollectionImpl struct {
	handlerFuncs []HandlerFunc
}

func newHandlerFuncCollection() HandlerFuncCollectionImpl {
	return HandlerFuncCollectionImpl{}
}

func (h *HandlerFuncCollectionImpl) AddHandlerFuncs(handlerFuncs []HandlerFunc) {
	h.handlerFuncs = append(h.handlerFuncs, handlerFuncs...)
}

func (h *HandlerFuncCollectionImpl) GetHandlerFunc(idx int) HandlerFunc {
	return h.handlerFuncs[idx]
}

func (h *HandlerFuncCollectionImpl) Len() int {
	return len(h.handlerFuncs)
}

type PatternHandlerFuncCollection struct {
	HandlerFuncCollectionImpl
	pattern string
}

func newPatternHandlerFuncCollection(pattern string) PatternHandlerFuncCollection {
	return PatternHandlerFuncCollection{
		HandlerFuncCollectionImpl: newHandlerFuncCollection(),
		pattern:                   pattern,
	}
}

type AllHandlerFuncCollection struct {
	HandlerFuncCollectionImpl
}

func newAllHandlerCollection() AllHandlerFuncCollection {
	return AllHandlerFuncCollection{newHandlerFuncCollection()}
}
