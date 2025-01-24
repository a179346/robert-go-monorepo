package middlewares

import (
	"net/http"

	"github.com/a179346/robert-go-monorepo/internal/post_board/shared/auth_jwt"
	"github.com/a179346/robert-go-monorepo/internal/post_board/shared/authed_context"
	"github.com/a179346/robert-go-monorepo/pkg/gohf_extended"
	"github.com/gohf-http/gohf/v6"
	"github.com/ztrue/tracerr"
)

func AuthedMiddleware(c *gohf.Context) gohf.Response {
	token := c.Req.GetHeader("auth_token")
	if token == "" {
		return gohf_extended.NewErrorResponse(
			http.StatusUnauthorized,
			"Unauthorized",
			tracerr.New("token is required"),
			false,
		)
	}

	claims, err := auth_jwt.Parse(token)
	if err != nil {
		return gohf_extended.NewErrorResponse(
			http.StatusUnauthorized,
			"Unauthorized",
			tracerr.Errorf("jwt parse error: %w", err),
			false,
		)
	}

	ctx := c.Req.Context()
	c.Req.SetContext(authed_context.WithValue(ctx, claims.ID))
	return c.Next()
}
