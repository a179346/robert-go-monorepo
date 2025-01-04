package roberthttp

import (
	"net/http"
)

func getHTTPHandleFunc(router Router, handlerFuncCollection HandlerFuncCollection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := newResponse(w, router.options.Response)
		c := newContext(res, newRequest(r))
		var handle func(idx int)
		handle = func(idx int) {
			if idx == handlerFuncCollection.Len() {
				return
			}
			c.Next = func() { handle(idx + 1) }

			handlerFuncCollection.GetHandlerFunc(idx)(c)
		}

		handle(0)
	}
}
