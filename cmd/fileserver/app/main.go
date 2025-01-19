package main

import (
	"context"
	"net/http"
	"os"
	"time"

	fileserver_applogger "github.com/a179346/robert-go-monorepo/internal/fileserver/applogger"
	_ "github.com/a179346/robert-go-monorepo/internal/fileserver/config"
	fileserver_config "github.com/a179346/robert-go-monorepo/internal/fileserver/config"
	fileserver_server "github.com/a179346/robert-go-monorepo/internal/fileserver/server"
	filestore_use_case "github.com/a179346/robert-go-monorepo/internal/fileserver/use_cases/filestore"
	"github.com/a179346/robert-go-monorepo/pkg/gohf_extended"
	"github.com/a179346/robert-go-monorepo/pkg/graceful_shutdown"
	"github.com/a179346/robert-go-monorepo/pkg/logger"
	"github.com/ztrue/tracerr"
)

func main() {
	tracerr.DefaultCap = 8
	gohf_extended.SetReponseErrorDetail(fileserver_config.GetDebugConfig().ResponseErrorDetail)
	appLogger := fileserver_applogger.GetFlushLogger()
	if appLogger != nil {
		gohf_extended.SetLogger(appLogger)
	}

	server := fileserver_server.New(
		fileserver_server.Options{
			FileStoreUseCase: filestore_use_case.New(),
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
