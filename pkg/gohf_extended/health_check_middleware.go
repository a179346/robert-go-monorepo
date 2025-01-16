package gohf_extended

import (
	"net/http"

	"github.com/gohf-http/gohf/v6"
	"github.com/gohf-http/gohf/v6/response"
)

func HealthzMiddleware(c *gohf.Context) gohf.Response {
	return response.Status(http.StatusOK)
}
