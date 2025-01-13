package middlewares

import (
	"errors"
	"net/http"

	"github.com/a179346/robert-go-monorepo/internal/post_board/shared/auth_jwt"
	"github.com/a179346/robert-go-monorepo/internal/post_board/shared/authed_context"
	"github.com/gohf-http/gohf/v6"
	"github.com/gohf-http/gohf/v6/response"
)

func AuthedMiddleware(c *gohf.Context) gohf.Response {
	token := c.Req.GetHeader("auth_token")
	if token == "" {
		return response.Error(http.StatusUnauthorized, errors.New("Unauthorized"))
	}

	claims, err := auth_jwt.Parse(token)
	if err != nil {
		return response.Error(http.StatusUnauthorized, errors.New("Unauthorized"))
	}

	ctx := c.Req.Context()
	c.Req.SetContext(authed_context.WithValue(ctx, claims.ID))
	return c.Next()
}
