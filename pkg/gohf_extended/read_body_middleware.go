package gohf_extended

import (
	"context"
	"io"
	"net/http"

	"github.com/gohf-http/gohf/v6"
	"github.com/ztrue/tracerr"
)

func ReadBodyMiddleware(c *gohf.Context) gohf.Response {
	defer c.Req.GetBody().Close()
	data, err := io.ReadAll(c.Req.GetBody())
	if err != nil {
		return NewErrorResponse(
			http.StatusInternalServerError,
			"Something went wrong",
			tracerr.Errorf("read body failed: %w", err),
		)
	}

	ctx := c.Req.Context()
	c.Req.SetContext(withBody(ctx, data))
	return c.Next()
}

type bodyContextKey struct{}

func withBody(ctx context.Context, data []byte) context.Context {
	return context.WithValue(ctx, bodyContextKey{}, data)
}

func BodyValue(ctx context.Context) ([]byte, bool) {
	value, ok := ctx.Value(bodyContextKey{}).([]byte)
	return value, ok
}
