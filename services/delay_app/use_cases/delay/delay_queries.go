package delay_use_case

import (
	"context"
	"time"

	"github.com/ztrue/tracerr"
)

type delayQueries struct{}

func newDelayQueries() delayQueries {
	return delayQueries{}
}

func (delayQueries delayQueries) getResult(ctx context.Context, ms int, d string) (string, error) {
	select {
	case <-time.After(time.Duration(ms) * time.Millisecond):
		return d, nil

	case <-ctx.Done():
		return "", tracerr.Errorf("get result error: %w", ctx.Err())
	}
}
