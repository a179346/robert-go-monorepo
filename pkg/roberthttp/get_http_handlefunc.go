package roberthttp

import (
	"net/http"
)

func getHTTPHandleFunc(router Router, handlerGroup *handlerGroup) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := newResponse(w)
		req := newRequest(r, router.fullprefix)
		c := newContext(res, req)

		var handle func(idx int) HttpResponse
		handle = func(idx int) HttpResponse {
			if idx == handlerGroup.len() {
				return nil
			}
			c.Next = func() HttpResponse { return handle(idx + 1) }

			handlerFuncWithPrefix := handlerGroup.getHandlerFuncWithPrefix(idx)

			originalPrefix := req.getCurrPrefix()
			req.setCurrPrefix(handlerFuncWithPrefix.prefix)
			httpResponse := handlerFuncWithPrefix.f(c)
			req.setCurrPrefix(originalPrefix)

			return httpResponse
		}

		if httpResponse := handle(0); httpResponse != nil {
			req.setCurrPrefix("")
			httpResponse.Send(c.Res, c.Req)
		}
	}
}
