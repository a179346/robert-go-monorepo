package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/a179346/robert-go-monorepo/pkg/console"
	"github.com/a179346/robert-go-monorepo/pkg/gohf_extended"
	"github.com/a179346/robert-go-monorepo/pkg/graceful"
	fileserver_apilogger "github.com/a179346/robert-go-monorepo/services/fileserver/apilogger"
	fileserver_config "github.com/a179346/robert-go-monorepo/services/fileserver/config"
	fileserver_server "github.com/a179346/robert-go-monorepo/services/fileserver/server"
	filestore_use_case "github.com/a179346/robert-go-monorepo/services/fileserver/use_cases/filestore"
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
	gohf_extended.SetReponseErrorDetail(fileserver_config.GetDebugConfig().ResponseErrorDetail)

	apiLogger := fileserver_apilogger.GetApiLogger()
	if apiLogger != nil {
		defer func() {
			console.Info("Shutting down app logger...")
			time.Sleep(2 * time.Second)
			apiLogger.Close()
		}()
	}

	server := fileserver_server.New(
		fileserver_server.Options{
			FileStoreUseCase: filestore_use_case.New(),
			ApiLogger:        apiLogger,
		},
	)

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
		defer cancel()
		console.Info("Shutting down server...")
		if err := server.Shutdown(ctx); err != nil {
			console.Errorf("Error shutting down server: %v", err)
		}
	}()

	serverListenErrCh := make(chan error)
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
