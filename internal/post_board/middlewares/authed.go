package middlewares

import (
	"context"
	"errors"
	"net/http"

	"github.com/a179346/robert-go-monorepo/internal/post_board/providers/jwt_provider"
	"github.com/a179346/robert-go-monorepo/pkg/gohf"
	"github.com/a179346/robert-go-monorepo/pkg/gohf/gohf_responses"
)

func AuthedMiddleware(jwtProvider jwt_provider.JwtProvider) gohf.HandlerFunc {
	unauthorizedResponse := gohf_responses.NewErrorResponse(http.StatusUnauthorized, errors.New("Unauthorized"))

	return func(c *gohf.Context) gohf.Response {
		token := c.Req.GetHeader("auth_token")
		if token == "" {
			return unauthorizedResponse
		}

		claims, err := jwtProvider.Parse(token)
		if err != nil {
			return unauthorizedResponse
		}

		ctx := c.Req.Context()
		c.Req.SetContext(withAuthed(ctx, claims.ID))
		return c.Next()
	}
}

type authedContextKey struct{}

func withAuthed(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, authedContextKey{}, userID)
}

func AuthedContextValue(ctx context.Context) (string, bool) {
	value, ok := ctx.Value(authedContextKey{}).(string)
	return value, ok
}
