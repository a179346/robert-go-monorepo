package main

import (
	"context"
	"log"
	"net/http"
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
)

func main() {
	db, err := dbhelper.Open()
	if err != nil {
		log.Fatalf("opendb.Open error: %v", err)
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
