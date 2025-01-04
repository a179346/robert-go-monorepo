package roberthttp

import (
	"net/http"

	"github.com/a179346/robert-go-monorepo/packages/set"
)

type Router struct {
	handlers []HandlerFuncCollection
	groups   []*RouterGroup
	patterns *set.Set[string]
	options  *RouterOptions
}

func New(options *RouterOptions) Router {
	options = defaultRouterOptions(options)

	return Router{
		patterns: set.New[string](),
		options:  options,
	}
}

func (r *Router) Use(handlerFuncs ...HandlerFunc) {
	handler := newAllHandlerCollection()
	handler.AddHandlerFuncs(handlerFuncs)

	r.handlers = append(r.handlers, &handler)
	for _, group := range r.groups {
		group.handlers = append(group.handlers, &handler)
	}
}

func (r *Router) Handle(pattern string, handlerFuncs ...HandlerFunc) {
	handler := newPatternHandlerFuncCollection(pattern)
	handler.AddHandlerFuncs(handlerFuncs)
	r.handlers = append(r.handlers, &handler)

	r.patterns.Add(pattern)
}

func (r *Router) Group(prefix string) *RouterGroup {
	group := newGroup(prefix, r.options)
	for _, handler := range r.handlers {
		if h, ok := handler.(*AllHandlerFuncCollection); ok {
			group.handlers = append(group.handlers, h)
		}
	}

	r.groups = append(r.groups, group)

	return group
}

func (r Router) CreateHttpHandler() http.Handler {
	mux := http.NewServeMux()

	for pattern := range r.patterns.All() {
		mergedHandler := newHandlerFuncCollection()
		for _, handler := range r.handlers {
			if h, ok := handler.(*PatternHandlerFuncCollection); ok && h.pattern == pattern {
				mergedHandler.AddHandlerFuncs(h.handlerFuncs)
				continue
			}
			if h, ok := handler.(*AllHandlerFuncCollection); ok {
				mergedHandler.AddHandlerFuncs(h.handlerFuncs)
				continue
			}
		}

		mux.HandleFunc(pattern, getHTTPHandleFunc(r, &mergedHandler))
	}

	for _, group := range r.groups {
		prefix := group.prefix
		mux.Handle(prefix+"/", http.StripPrefix(prefix, group.CreateHttpHandler()))
	}

	return mux
}
