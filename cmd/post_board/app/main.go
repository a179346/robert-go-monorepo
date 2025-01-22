package main

import (
	"context"
	"fmt"
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

	gohf_extended.SetAppId("post_board")
	gohf_extended.SetReponseErrorDetail(post_board_config.GetDebugConfig().ResponseErrorDetail)
	appLogger, err := post_board_applogger.GetAppLogger()
	if err != nil {
		return fmt.Errorf("GetAppLogger error: %w", err)
	}
	if appLogger != nil {
		defer func() {
			console.Info("Shutting down app logger...")
			time.Sleep(2 * time.Second)
			appLogger.Close()
		}()
		gohf_extended.SetAppLogger(appLogger)
	}

	db, err := dbhelper.Open()
	if err != nil {
		return fmt.Errorf("opendb.Open error: %w", err)
	}
	defer func() {
		console.Info("Shutting down db...")
		db.Close()
	}()
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
