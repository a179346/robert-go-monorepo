package calculator

import (
	"context"
	"errors"

	"github.com/ztrue/tracerr"
)

type Service interface {
	divide(ctx context.Context, dividend, divisor float64) (float64, error)
}

type service struct{}

func NewService() *service {
	return &service{}
}

var errDividedByZero = errors.New("divide by 0")

func (service service) divide(ctx context.Context, dividend, divisor float64) (float64, error) {
	if divisor == 0 {
		return 0, tracerr.Wrap(errDividedByZero)
	}

	return dividend / divisor, nil
}
