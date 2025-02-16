package calculator

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	divide endpoint.Endpoint
}

func NewEndpoints(s Service) Endpoints {
	return Endpoints{
		divide: makeDivideEndpoint(s),
	}
}

type divideReq struct {
	Dividend float64 `json:"dividend"`
	Divisor  float64 `json:"divisor" validate:"required"`
}

type divideResp struct {
	Value float64 `json:"value"`
}

func makeDivideEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(divideReq)
		if !ok {
			return nil, errors.New("request should be of type divideReq")
		}

		value, err := s.divide(ctx, req.Dividend, req.Divisor)
		if err != nil {
			return divideResp{value}, err
		}
		return divideResp{value}, nil
	}
}
