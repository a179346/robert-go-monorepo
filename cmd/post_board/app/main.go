package main

import (
	"context"
	"net/http"
	"os"
	"time"

	post_board_applogger "github.com/a179346/robert-go-monorepo/internal/post_board/applogger"
	post_board_config "github.com/a179346/robert-go-monorepo/internal/post_board/config"
	"github.com/a179346/robert-go-monorepo/internal/post_board/database/dbhelper"
	"github.com/a179346/robert-go-monorepo/internal/post_board/providers/post_provider"
	"github.com/a179346/robert-go-monorepo/internal/post_board/providers/user_provider"
	post_board_server "github.com/a179346/robert-go-monorepo/internal/post_board/server"
	auth_use_case "github.com/a179346/robert-go-monorepo/internal/post_board/use_cases/auth"
	post_use_case "github.com/a179346/robert-go-monorepo/internal/post_board/use_cases/post"
	user_use_case "github.com/a179346/robert-go-monorepo/internal/post_board/use_cases/user"
	"github.com/a179346/robert-go-monorepo/pkg/gohf_extended"
	"github.com/a179346/robert-go-monorepo/pkg/graceful_shutdown"
	"github.com/a179346/robert-go-monorepo/pkg/logger"
	"github.com/ztrue/tracerr"
)

func main() {
	tracerr.DefaultCap = 8
	gohf_extended.SetReponseErrorDetail(post_board_config.GetDebugConfig().ResponseErrorDetail)
	appLogger, err := post_board_applogger.GetRabbitMQLogger()
	if err != nil {
		logger.Errorf("GetRabbitMQLogger error: %v", err)
		os.Exit(1)
	}
	if appLogger != nil {
		gohf_extended.SetLogger(appLogger)
	}

	db, err := dbhelper.Open()
	if err != nil {
		logger.Errorf("opendb.Open error: %v", err)
		os.Exit(1)
	}
	db.SetMaxOpenConns(30)
	dbhelper.WaitFor(context.Background(), db)

	userProvider := user_provider.New(db)
	postProvider := post_provider.New(db)

	server := post_board_server.New(
		post_board_server.Options{
			AuthUseCase: auth_use_case.New(userProvider),
			UserUseCase: user_use_case.New(userProvider),
			PostUseCase: post_use_case.New(postProvider),
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

	logger.Info("Shutting down db...")
	db.Close()

	if appLogger != nil {
		logger.Info("Shutting down app logger...")
		time.Sleep(2 * time.Second)
		appLogger.Close(ctx)
	}

	logger.Info("Server shut down successfully")
}
