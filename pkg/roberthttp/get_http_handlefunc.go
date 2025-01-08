package roberthttp

import (
	"net/http"
)

func getHTTPHandleFunc(router *Router, handlers []*handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := newResponse(w)
		req := newRequest(r, router.fullPrefix)
		c := newContext(res, req)

		var handle func(idx int) HttpResponse
		handle = func(idx int) HttpResponse {
			if idx == len(handlers) {
				return nil
			}

			c.Next = func() HttpResponse { return handle(idx + 1) }

			handler := handlers[idx]

			originalPrefix := req.getCurrPrefix()
			req.setCurrPrefix(handler.owner.fullPrefix)
			httpResponse := handler.f(c)
			req.setCurrPrefix(originalPrefix)

			return httpResponse
		}

		if httpResponse := handle(0); httpResponse != nil {
			req.setCurrPrefix("")
			httpResponse.Send(c.Res, c.Req)
		}
	}
}
