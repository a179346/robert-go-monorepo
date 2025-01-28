package authed_context

import "context"

type authedContextKey struct{}

func WithValue(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, authedContextKey{}, userID)
}

func Value(ctx context.Context) (string, bool) {
	value, ok := ctx.Value(authedContextKey{}).(string)
	return value, ok
}
