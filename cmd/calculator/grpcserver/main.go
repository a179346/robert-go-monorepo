package main

import (
	"fmt"
	"net"
	"os"

	calculatorPb "github.com/a179346/robert-go-monorepo/pb/calculator"
	"github.com/a179346/robert-go-monorepo/pkg/console"
	"github.com/a179346/robert-go-monorepo/pkg/graceful"
	"github.com/a179346/robert-go-monorepo/services/calculator"
	"github.com/ztrue/tracerr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	if err := run(); err != nil {
		console.Errorf("Exit 1: %v", err)
		os.Exit(1)
	}
	console.Info("Exit 0")
}

func run() error {
	tracerr.DefaultCap = 8

	calculatorService := calculator.NewService()
	calculatorEndpoints := calculator.NewEndpoints(calculatorService)
	calculatorServer := calculator.NewGRPCServer(calculatorEndpoints)

	grpcListener, err := net.Listen("tcp", ":9085")
	if err != nil {
		return fmt.Errorf("net.Listen error: %w", err)
	}
	defer grpcListener.Close()

	grpcServer := grpc.NewServer()
	calculatorPb.RegisterCalculatorServer(grpcServer, calculatorServer)
	reflection.Register(grpcServer)
	defer grpcServer.GracefulStop()

	serverListenErrCh := make(chan error)
	go func() {
		if err := grpcServer.Serve(grpcListener); err != nil {
			serverListenErrCh <- fmt.Errorf("grpcServer.Serve error: %w", err)
		}
	}()

	select {
	case signal := <-graceful.ShutDown():
		console.Infof("Received signal: %v", signal)
		return nil

	case err := <-serverListenErrCh:
		return err
	}
}
