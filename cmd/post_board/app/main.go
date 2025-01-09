package main

import (
	"context"
	"log"
	"net/http"
	"time"

	post_board_config "github.com/a179346/robert-go-monorepo/internal/post_board/config"
	"github.com/a179346/robert-go-monorepo/internal/post_board/database/opendb"
	"github.com/a179346/robert-go-monorepo/internal/post_board/providers/user_provider"
	post_board_server "github.com/a179346/robert-go-monorepo/internal/post_board/server"
	auth_use_case "github.com/a179346/robert-go-monorepo/internal/post_board/use_caes/auth"
	"github.com/a179346/robert-go-monorepo/pkg/graceful_shutdown"
)

func main() {
	config := post_board_config.New()

	// TODO - check db connected
	db, err := opendb.Open(config.DB)
	if err != nil {
		log.Fatalf("opendb.Open error: %v", err)
	}

	userProvider := user_provider.New(db)

	server := post_board_server.New(
		config.Server,
		post_board_server.Options{
			AuthUseCase: auth_use_case.New(userProvider),
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
