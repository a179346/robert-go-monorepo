package main

import (
	"context"
	"log"
	"net/http"
	"time"

	post_board_config "github.com/a179346/robert-go-monorepo/internal/post_board/config"
	"github.com/a179346/robert-go-monorepo/internal/post_board/database/dbhelper"
	"github.com/a179346/robert-go-monorepo/internal/post_board/middlewares"
	"github.com/a179346/robert-go-monorepo/internal/post_board/providers/jwt_provider"
	"github.com/a179346/robert-go-monorepo/internal/post_board/providers/user_provider"
	post_board_server "github.com/a179346/robert-go-monorepo/internal/post_board/server"
	auth_use_case "github.com/a179346/robert-go-monorepo/internal/post_board/use_caes/auth"
	user_use_case "github.com/a179346/robert-go-monorepo/internal/post_board/use_caes/user"
	"github.com/a179346/robert-go-monorepo/pkg/graceful_shutdown"
)

func main() {
	config := post_board_config.New()

	db, err := dbhelper.Open(config.DB)
	if err != nil {
		log.Fatalf("opendb.Open error: %v", err)
	}
	db.SetMaxOpenConns(30)
	dbhelper.WaitFor(db)

	userProvider := user_provider.New(db)
	jwtProvider := jwt_provider.New(config.Jwt.Secret, config.Jwt.ExpireSeconds)

	server := post_board_server.New(
		config.Server,
		post_board_server.Options{
			AuthedMiddleware: middlewares.AuthedMiddleware(jwtProvider),
			AuthUseCase:      auth_use_case.New(userProvider, jwtProvider),
			UserUseCase:      user_use_case.New(userProvider),
		},
	)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	signal := <-graceful_shutdown.ShutDown()
	log.Printf("Received signal: %v", signal)
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Error shutting down server: %v", err)
	}
	db.Close()

	log.Println("Server shut down successfully")
}
