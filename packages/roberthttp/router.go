package roberthttp

import (
	"net/http"
)

type empty struct{}

type Router struct {
	handlers []Handler
	groups   []RouterGroup
	patterns map[string]empty
}

func New() Router {
	return Router{patterns: make(map[string]empty)}
}

func (r *Router) Use(handlerFuncs ...HandlerFunc) {
	handler := newAllHandler()
	handler.AddHandlerFuncs(handlerFuncs)
	r.handlers = append(r.handlers, &handler)
}

func (r *Router) Handle(pattern string, handlerFuncs ...HandlerFunc) {
	handler := newPatternHandler(pattern)
	handler.AddHandlerFuncs(handlerFuncs)
	r.handlers = append(r.handlers, &handler)

	r.patterns[pattern] = empty{}
}

func (r *Router) Group(prefix string) RouterGroup {
	group := NewGroup(prefix)
	r.groups = append(r.groups, group)
	return group
}

func (r *Router) AddGroup(group RouterGroup) {
	r.groups = append(r.groups, group)
}

func (r Router) CreateHttpHandler() http.Handler {
	mux := http.NewServeMux()

	for pattern := range r.patterns {
		mergedHandler := newHandler()
		for _, handler := range r.handlers {
			if h, ok := handler.(*PatternHandler); ok && h.pattern == pattern {
				mergedHandler.AddHandlerFuncs(h.handlerFuncs)
				continue
			}
			if h, ok := handler.(*AllHandler); ok {
				mergedHandler.AddHandlerFuncs(h.handlerFuncs)
				continue
			}
		}

		mux.Handle(pattern, mergedHandler)
	}

	for _, group := range r.groups {
		prefix := group.prefix
		mux.Handle(prefix+"/", http.StripPrefix(prefix, group.CreateHttpHandler()))
	}

	return mux
}
