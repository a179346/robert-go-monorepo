package main

import (
	"context"
	"net/http"
	"os"
	"time"

	_ "github.com/a179346/robert-go-monorepo/internal/post_board/config"
	"github.com/a179346/robert-go-monorepo/internal/post_board/database/dbhelper"
	"github.com/a179346/robert-go-monorepo/internal/post_board/providers/post_provider"
	"github.com/a179346/robert-go-monorepo/internal/post_board/providers/user_provider"
	post_board_server "github.com/a179346/robert-go-monorepo/internal/post_board/server"
	auth_use_case "github.com/a179346/robert-go-monorepo/internal/post_board/use_cases/auth"
	post_use_case "github.com/a179346/robert-go-monorepo/internal/post_board/use_cases/post"
	user_use_case "github.com/a179346/robert-go-monorepo/internal/post_board/use_cases/user"
	"github.com/a179346/robert-go-monorepo/pkg/graceful_shutdown"
	"github.com/a179346/robert-go-monorepo/pkg/logger"
)

func main() {
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
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Errorf("Error shutting down server: %v", err)
	}
	db.Close()

	logger.Info("Server shut down successfully")
}
