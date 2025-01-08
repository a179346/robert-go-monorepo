package roberthttp_extended

import (
	"errors"
	"net/http"

	"github.com/a179346/robert-go-monorepo/pkg/roberthttp"
	"github.com/a179346/robert-go-monorepo/pkg/roberthttp/roberthttp_response"
)

func NotFoundHandler(c *roberthttp.Context) roberthttp.HttpResponse {
	return roberthttp_response.NewErrorResponse(
		http.StatusNotFound,
		errors.New("Not found"),
	)
}
