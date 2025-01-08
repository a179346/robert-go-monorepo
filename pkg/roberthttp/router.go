package roberthttp

import (
	"net/http"

	"github.com/a179346/robert-go-monorepo/pkg/set"
)

type Router struct {
	prefix        string
	fullprefix    string
	handlerGroups []*handlerGroup
	subRouters    []*Router
	patterns      *set.Set[string]
}

func New() *Router {
	patterns := set.New[string]()
	patterns.Add("/")
	return &Router{
		patterns: patterns,
	}
}

func (r *Router) Use(handlerFuncs ...HandlerFunc) {
	handlerGroup := newHandlerGroup("", true)
	for _, f := range handlerFuncs {
		handlerFuncWithPrefix := newHandlerFuncWithPrefix(f, r.fullprefix)
		handlerGroup.addHandlerFuncWithPrefix(handlerFuncWithPrefix)
	}
	r.appendHandlerGroup(&handlerGroup)
}

func (r *Router) Handle(pattern string, handlerFuncs ...HandlerFunc) {
	handlerGroup := newHandlerGroup(pattern, false)
	for _, f := range handlerFuncs {
		handlerFuncWithPrefix := newHandlerFuncWithPrefix(f, r.fullprefix)
		handlerGroup.addHandlerFuncWithPrefix(handlerFuncWithPrefix)
	}

	r.appendHandlerGroup(&handlerGroup)

	r.patterns.Add(pattern)
}

func (r *Router) SubRouter(prefix string) *Router {
	subRouter := New()

	if len(prefix) > 0 && prefix[len(prefix)-1] == '/' {
		prefix = prefix[:len(prefix)-1]
	}
	subRouter.prefix = prefix
	subRouter.fullprefix = r.fullprefix + prefix

	for _, handlerGroup := range r.handlerGroups {
		if handlerGroup.all {
			subRouter.appendHandlerGroup(handlerGroup)
		}
	}

	r.subRouters = append(r.subRouters, subRouter)

	return subRouter
}

func (r Router) CreateHttpHandler() http.Handler {
	mux := http.NewServeMux()

	handlerGroupMap := r.getHandlerGroupMap()

	for _, handlerGroup := range handlerGroupMap {
		if handlerGroup.len() > 0 {
			mux.HandleFunc(handlerGroup.pattern, getHTTPHandleFunc(r, handlerGroup))
		}
	}

	for _, router := range r.subRouters {
		prefix := router.prefix
		mux.Handle(prefix+"/", http.StripPrefix(prefix, router.CreateHttpHandler()))
	}

	return mux
}

func (r Router) getHandlerGroupMap() map[string]*handlerGroup {
	handlerGroupMap := make(map[string]*handlerGroup)

	for pattern := range r.patterns.All() {
		mergedHandlerGroup := newHandlerGroup(pattern, false)
		handlerGroupMap[pattern] = &mergedHandlerGroup
	}

	for _, handlerGroup := range r.handlerGroups {
		if handlerGroup.all {
			for _, mergedHandlerGroup := range handlerGroupMap {
				mergedHandlerGroup.copyHandlerFuncWithPrefixsFrom(handlerGroup)
			}
		} else {
			mergedHandlerGroup := handlerGroupMap[handlerGroup.pattern]
			mergedHandlerGroup.copyHandlerFuncWithPrefixsFrom(handlerGroup)
		}
	}

	return handlerGroupMap
}

func (r *Router) appendHandlerGroup(handlerGroup *handlerGroup) {
	r.handlerGroups = append(r.handlerGroups, handlerGroup)
	if handlerGroup.all {
		for _, router := range r.subRouters {
			router.appendHandlerGroup(handlerGroup)
		}
	}
}
