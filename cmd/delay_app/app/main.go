package main

import (
	"context"
	"log"
	"net/http"
	"time"

	_ "github.com/a179346/robert-go-monorepo/internal/delay_app/config"
	delay_app_server "github.com/a179346/robert-go-monorepo/internal/delay_app/server"
	delay_use_case "github.com/a179346/robert-go-monorepo/internal/delay_app/use_caes/delay"
	"github.com/a179346/robert-go-monorepo/pkg/graceful_shutdown"
)

func main() {
	server := delay_app_server.New(
		delay_app_server.Options{
			DelayUseCase: delay_use_case.New(),
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

	log.Println("Server shut down successfully")
}
