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
	post_board_apilogger "github.com/a179346/robert-go-monorepo/services/post_board/apilogger"
	post_board_config "github.com/a179346/robert-go-monorepo/services/post_board/config"
	"github.com/a179346/robert-go-monorepo/services/post_board/database/dbhelper"
	"github.com/a179346/robert-go-monorepo/services/post_board/providers/post_provider"
	"github.com/a179346/robert-go-monorepo/services/post_board/providers/user_provider"
	post_board_server "github.com/a179346/robert-go-monorepo/services/post_board/server"
	auth_use_case "github.com/a179346/robert-go-monorepo/services/post_board/use_cases/auth"
	post_use_case "github.com/a179346/robert-go-monorepo/services/post_board/use_cases/post"
	user_use_case "github.com/a179346/robert-go-monorepo/services/post_board/use_cases/user"
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
	gohf_extended.SetReponseErrorDetail(post_board_config.GetDebugConfig().ResponseErrorDetail)

	apiLogger, err := post_board_apilogger.GetApiLogger()
	if err != nil {
		return fmt.Errorf("GetApiLogger error: %w", err)
	}
	if apiLogger != nil {
		defer func() {
			console.Info("Shutting down app logger...")
			time.Sleep(2 * time.Second)
			apiLogger.Close()
		}()
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
			ApiLogger:   apiLogger,
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
