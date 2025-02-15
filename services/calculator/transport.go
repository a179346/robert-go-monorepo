package calculator

import (
	"context"
	"errors"

	calculatorPb "github.com/a179346/robert-go-monorepo/pb/calculator"
	gt "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type gRPCServer struct {
	divide gt.Handler
	calculatorPb.UnimplementedCalculatorServer
}

func NewGRPCServer(endpoints Endpoints) calculatorPb.CalculatorServer {
	return &gRPCServer{
		divide: gt.NewServer(
			endpoints.divide,
			decodeDivideRequest,
			encodeDivideResponse,
		),
	}
}

func (s gRPCServer) Divide(ctx context.Context, req *calculatorPb.DivideRequest) (*calculatorPb.DivideResponse, error) {
	_, resp, err := s.divide.ServeGRPC(ctx, req)
	switch {
	case err == nil:
		return resp.(*calculatorPb.DivideResponse), nil

	case errors.Is(err, errDividedByZero):
		return nil, status.Error(codes.InvalidArgument, err.Error())

	default:
		return nil, status.Error(codes.Unknown, err.Error())
	}
}

func decodeDivideRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*calculatorPb.DivideRequest)
	if !ok {
		return nil, errors.New("invalid request body")
	}

	return divideReq{Dividend: req.Dividend, Divisor: req.Divisor}, nil
}

func encodeDivideResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(divideResp)
	if !ok {
		return nil, errors.New("invalid response body")
	}

	return &calculatorPb.DivideResponse{Value: resp.Value}, nil
}
