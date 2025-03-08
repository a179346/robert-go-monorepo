package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/a179346/robert-go-monorepo/pkg/console"
	"github.com/a179346/robert-go-monorepo/pkg/graceful"
	"github.com/a179346/robert-go-monorepo/services/calculator"
	"github.com/ztrue/tracerr"
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
	httpServer := calculator.NewHttpServer(9084, calculatorEndpoints)

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
		defer cancel()
		console.Info("Shutting down server...")
		if err := httpServer.Shutdown(ctx); err != nil {
			console.Errorf("Error shutting down server: %v", err)
		}
	}()

	serverListenErrCh := make(chan error)
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverListenErrCh <- fmt.Errorf("Error starting server: %w", err)
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
