package roberthttp

import (
	"net/http"
)

func getHTTPHandleFunc(_ Router, handlerFuncCollection *HandlerFuncCollection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := newResponse(w)
		c := newContext(res, newRequest(r))
		var handle func(idx int) HttpResponse
		handle = func(idx int) HttpResponse {
			if idx == handlerFuncCollection.Len() {
				return nil
			}
			c.Next = func() HttpResponse { return handle(idx + 1) }

			return handlerFuncCollection.GetHandlerFunc(idx)(c)
		}

		if httpResponse := handle(0); httpResponse != nil {
			httpResponse.Send(c.Res, c.Req)
		}
	}
}
