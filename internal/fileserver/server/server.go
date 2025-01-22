package fileserver_server

import (
	"context"
	"fmt"
	"net/http"

	fileserver_config "github.com/a179346/robert-go-monorepo/internal/fileserver/config"
	filestore_use_case "github.com/a179346/robert-go-monorepo/internal/fileserver/use_cases/filestore"
	"github.com/a179346/robert-go-monorepo/pkg/console"
	"github.com/a179346/robert-go-monorepo/pkg/gohf_extended"
	"github.com/gohf-http/gohf/v6"
	"github.com/rs/cors"
)

type Options struct {
	FileStoreUseCase filestore_use_case.FileStoreUseCase
}

type Server struct {
	httpserver *http.Server
}

func New(options Options) *Server {
	router := gohf.New()

	router.Use(gohf.MaxBytesMiddleware(5 * 1024 * 1024))
	router.Use(gohf_extended.RequestIdMiddleware)
	router.Use(gohf_extended.ReadBodyMiddleware)

	router.GET("/healthz", gohf_extended.HealthzHandler)

	options.FileStoreUseCase.AppendHandler(router.SubRouter("/filestore"))

	router.Use(gohf_extended.NotFoundHandler)

	mux := router.CreateServeMux()
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", fileserver_config.GetServerConfig().Port),
		Handler: cors.AllowAll().Handler(mux),
	}

	return &Server{httpserver: server}
}

func (s *Server) ListenAndServe() error {
	console.Infof("Starting server on \"%s\"", s.httpserver.Addr)
	return s.httpserver.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpserver.Shutdown(ctx)
}
