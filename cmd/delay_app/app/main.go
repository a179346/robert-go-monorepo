package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	delay_app_applogger "github.com/a179346/robert-go-monorepo/internal/delay_app/applogger"
	_ "github.com/a179346/robert-go-monorepo/internal/delay_app/config"
	delay_app_config "github.com/a179346/robert-go-monorepo/internal/delay_app/config"
	delay_app_server "github.com/a179346/robert-go-monorepo/internal/delay_app/server"
	delay_use_case "github.com/a179346/robert-go-monorepo/internal/delay_app/use_cases/delay"
	"github.com/a179346/robert-go-monorepo/pkg/console"
	"github.com/a179346/robert-go-monorepo/pkg/gohf_extended"
	"github.com/a179346/robert-go-monorepo/pkg/graceful_shutdown"
	"github.com/ztrue/tracerr"
)

func main() {
	if err := run(); err != nil {
		console.Errorf("%v", err)
		os.Exit(1)
	}
}

func run() error {
	tracerr.DefaultCap = 8

	gohf_extended.SetAppId("delay_app")
	gohf_extended.SetReponseErrorDetail(delay_app_config.GetDebugConfig().ResponseErrorDetail)
	appLogger := delay_app_applogger.GetAppLogger()
	if appLogger != nil {
		defer func() {
			console.Info("Shutting down app logger...")
			time.Sleep(2 * time.Second)
			appLogger.Close()
		}()
		gohf_extended.SetAppLogger(appLogger)
	}

	server := delay_app_server.New(
		delay_app_server.Options{
			DelayUseCase: delay_use_case.New(),
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
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverListenErrCh <- fmt.Errorf("Error starting server: %w", err)
		}
	}()

	select {
	case signal := <-graceful_shutdown.ShutDown():
		console.Infof("Received signal: %v", signal)
		return nil

	case err := <-serverListenErrCh:
		return err
	}
}
