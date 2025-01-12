package gohf_extended

import (
	"errors"
	"net/http"

	"github.com/gohf-http/gohf/v4"
	"github.com/gohf-http/gohf/v4/gohf_responses"
)

func NotFoundHandler(c *gohf.Context) gohf.Response {
	return gohf_responses.NewErrorResponse(
		http.StatusNotFound,
		errors.New("Not found"),
	)
}
