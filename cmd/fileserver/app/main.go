package main

import (
	"context"
	"log"
	"net/http"
	"time"

	fileserver_config "github.com/a179346/robert-go-monorepo/internal/fileserver/config"
	fileserver_server "github.com/a179346/robert-go-monorepo/internal/fileserver/server"
	filestore_use_case "github.com/a179346/robert-go-monorepo/internal/fileserver/use_caes/filestore"
	"github.com/a179346/robert-go-monorepo/pkg/graceful_shutdown"
)

func main() {
	config := fileserver_config.New()

	server := fileserver_server.New(
		config.Server,
		fileserver_server.Options{
			FileStoreUseCase: filestore_use_case.New(config.Storage),
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
