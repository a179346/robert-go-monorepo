package gohf_extended

import (
	"context"

	"github.com/gohf-http/gohf/v6"
	"github.com/google/uuid"
)

func RequestIdMiddleware(c *gohf.Context) gohf.Response {
	ctx := c.Req.Context()
	requestId := uuid.New()
	c.Req.SetContext(withId(ctx, requestId))
	return c.Next()
}

type requestIdContextKey struct{}

func withId(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, requestIdContextKey{}, id)
}

func IdValue(ctx context.Context) (uuid.UUID, bool) {
	value, ok := ctx.Value(requestIdContextKey{}).(uuid.UUID)
	return value, ok
}
