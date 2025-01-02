package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	delay_app_config "github.com/a179346/robert-go-monorepo/apps/delay_app/config"
	delay_app_server "github.com/a179346/robert-go-monorepo/apps/delay_app/server"
	delay_use_case "github.com/a179346/robert-go-monorepo/apps/delay_app/use_caes/delay"
	"github.com/a179346/robert-go-monorepo/packages/graceful_shutdown"
)

func main() {
	config := delay_app_config.New()

	server := delay_app_server.New(
		config.Server,
		delay_app_server.Options{
			DelayUseCase: delay_use_case.New(),
		},
	)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	graceful_shutdown.OnShutdown(func(sig os.Signal) {
		log.Printf("Received signal: %v", sig)
		log.Println("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down server: %v", err)
		}

		log.Println("Server shut down successfully")
	})
}
