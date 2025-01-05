package roberthttp

import (
	"net/http"

	"github.com/a179346/robert-go-monorepo/pkg/set"
)

type Router struct {
	prefix     string
	handlers   []*HandlerFuncCollection
	subRouters []*Router
	patterns   *set.Set[string]
	options    *RouterOptions
}

func New(options *RouterOptions) *Router {
	options = defaultRouterOptions(options)

	return &Router{
		patterns: set.New[string](),
		options:  options,
	}
}

func (r *Router) Use(handlerFuncs ...HandlerFunc) {
	handler := newHandlerFuncCollection("", true)
	handler.AddHandlerFuncs(handlerFuncs)
	r.handlers = append(r.handlers, &handler)
	for _, router := range r.subRouters {
		router.Use(handlerFuncs...)
	}
}

func (r *Router) Handle(pattern string, handlerFuncs ...HandlerFunc) {
	handler := newHandlerFuncCollection(pattern, false)
	handler.AddHandlerFuncs(handlerFuncs)
	r.handlers = append(r.handlers, &handler)

	r.patterns.Add(pattern)
}

func (r *Router) SubRouter(prefix string) *Router {
	subRouter := New(r.options)

	if len(prefix) > 0 && prefix[len(prefix)-1] == '/' {
		prefix = prefix[:len(prefix)-1]
	}
	subRouter.prefix = prefix

	for _, handler := range r.handlers {
		if handler.all {
			subRouter.handlers = append(subRouter.handlers, handler)
		}
	}

	r.subRouters = append(r.subRouters, subRouter)

	return subRouter
}

func (r Router) CreateHttpHandler() http.Handler {
	mux := http.NewServeMux()

	r.appendHandlersToServerMux(mux)

	for _, router := range r.subRouters {
		prefix := router.prefix
		mux.Handle(prefix+"/", http.StripPrefix(prefix, router.CreateHttpHandler()))
	}

	return mux
}

func (r Router) appendHandlersToServerMux(mux *http.ServeMux) {
	mergedHandlerMap := make(map[string]*HandlerFuncCollection)

	for pattern := range r.patterns.All() {
		mergedHandler := newHandlerFuncCollection(pattern, false)
		mergedHandlerMap[pattern] = &mergedHandler
	}

	for _, handler := range r.handlers {
		if handler.all {
			for _, mergedHandler := range mergedHandlerMap {
				mergedHandler.AddHandlerFuncs(handler.handlerFuncs)
			}
		} else {
			mergedHandlerMap[handler.pattern].AddHandlerFuncs(handler.handlerFuncs)
		}
	}

	for _, mergedHandler := range mergedHandlerMap {
		mux.HandleFunc(mergedHandler.pattern, getHTTPHandleFunc(r, mergedHandler))
	}
}
