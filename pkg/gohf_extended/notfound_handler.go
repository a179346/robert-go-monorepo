package gohf_extended

import (
	"errors"
	"net/http"

	"github.com/gohf-http/gohf/v6"
	"github.com/gohf-http/gohf/v6/response"
)

func NotFoundHandler(c *gohf.Context) gohf.Response {
	return response.Error(
		http.StatusNotFound,
		errors.New("Not found"),
	)
}
