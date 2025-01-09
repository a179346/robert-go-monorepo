package gohf_extended

import (
	"errors"
	"net/http"

	"github.com/a179346/robert-go-monorepo/pkg/gohf"
	"github.com/a179346/robert-go-monorepo/pkg/gohf/gohf_responses"
)

func NotFoundHandler(c *gohf.Context) gohf.Response {
	return gohf_responses.NewErrorResponse(
		http.StatusNotFound,
		errors.New("Not found"),
	)
}
