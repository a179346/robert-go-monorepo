package roberthttp

import (
	"net/http"
)

type Router struct {
	prefix            string
	fullPrefix        string
	handlerRepository *handlerRepository
	parentRouter      *Router
	subRouters        []*Router
}

func New() *Router {
	return &Router{
		handlerRepository: newHandlerRepository(),
	}
}

func (r *Router) Use(handlerFuncs ...HandlerFunc) {
	for _, f := range handlerFuncs {
		r.handlerRepository.addHandler(f, r, "", true)
	}
}

func (r *Router) Handle(pattern string, handlerFuncs ...HandlerFunc) {
	for _, f := range handlerFuncs {
		r.handlerRepository.addHandler(f, r, pattern, false)
	}
}

func (r *Router) SubRouter(prefix string) *Router {
	if len(prefix) > 0 && prefix[len(prefix)-1] == '/' {
		prefix = prefix[:len(prefix)-1]
	}

	subRouter := &Router{
		prefix:            prefix,
		fullPrefix:        r.fullPrefix + prefix,
		handlerRepository: r.handlerRepository,
		parentRouter:      r,
	}

	r.subRouters = append(r.subRouters, subRouter)

	return subRouter
}

func (r *Router) CreateHttpHandler() http.Handler {
	mux := http.NewServeMux()

	handlersMap := r.getHandlersMap()

	for pattern, handlers := range handlersMap {
		if len(handlers) > 0 {
			mux.HandleFunc(pattern, getHTTPHandleFunc(r.fullPrefix, handlers))
		}
	}

	for _, router := range r.subRouters {
		prefix := router.prefix
		mux.Handle(prefix+"/", http.StripPrefix(prefix, router.CreateHttpHandler()))
	}

	return mux
}

func (r *Router) getHandlersMap() map[string][]*handler {
	handlersMap := make(map[string][]*handler)
	handlersMap["/"] = make([]*handler, 0)

	for _, h := range r.handlerRepository.getHandlers() {
		if !h.all && h.owner == r {
			if _, ok := handlersMap[h.pattern]; !ok {
				handlersMap[h.pattern] = make([]*handler, 0)
			}
		}
	}

	for _, handler := range r.handlerRepository.getHandlers() {
		if handler.all && (handler.owner == r || r.isAncestor(handler.owner)) {
			for pattern := range handlersMap {
				handlersMap[pattern] = append(handlersMap[pattern], handler)
			}
		} else if !handler.all && handler.owner == r {
			pattern := handler.pattern
			handlersMap[pattern] = append(handlersMap[pattern], handler)
		}
	}

	return handlersMap
}

func (r *Router) isAncestor(router *Router) bool {
	parentRouter := r.parentRouter
	if parentRouter == router {
		return true
	}
	return parentRouter != nil && parentRouter.isAncestor(router)
}
