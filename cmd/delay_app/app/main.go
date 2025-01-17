package main

import (
	"context"
	"net/http"
	"os"
	"time"

	delay_app_applogger "github.com/a179346/robert-go-monorepo/internal/delay_app/applogger"
	_ "github.com/a179346/robert-go-monorepo/internal/delay_app/config"
	delay_app_server "github.com/a179346/robert-go-monorepo/internal/delay_app/server"
	delay_use_case "github.com/a179346/robert-go-monorepo/internal/delay_app/use_cases/delay"
	"github.com/a179346/robert-go-monorepo/pkg/gohf_extended"
	"github.com/a179346/robert-go-monorepo/pkg/graceful_shutdown"
	"github.com/a179346/robert-go-monorepo/pkg/logger"
)

func main() {
	appLogger := delay_app_applogger.GetFlushLogger()
	if appLogger != nil {
		gohf_extended.SetLogger(appLogger)
	}

	server := delay_app_server.New(
		delay_app_server.Options{
			DelayUseCase: delay_use_case.New(),
		},
	)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf("Error starting server: %v", err)
			os.Exit(1)
		}
	}()

	signal := <-graceful_shutdown.ShutDown()
	logger.Infof("Received signal: %v", signal)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	logger.Info("Shutting down server...")
	if err := server.Shutdown(ctx); err != nil {
		logger.Errorf("Error shutting down server: %v", err)
	}

	if appLogger != nil {
		logger.Info("Shutting down app logger...")
		time.Sleep(2 * time.Second)
		appLogger.Close(ctx)
	}

	logger.Info("Server shut down successfully")
}
