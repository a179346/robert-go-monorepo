package fileserver_server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	fileserver_config "github.com/a179346/robert-go-monorepo/internal/fileserver/config"
	filestore_use_case "github.com/a179346/robert-go-monorepo/internal/fileserver/use_caes/filestore"
	"github.com/a179346/robert-go-monorepo/pkg/roberthttp"
	"github.com/a179346/robert-go-monorepo/pkg/roberthttp_extended"
)

type Options struct {
	FileStoreUseCase filestore_use_case.FileStoreUseCase
}

type Server struct {
	httpserver *http.Server
}

func New(config fileserver_config.ServerConfig, options Options) *Server {
	router := roberthttp.New(roberthttp_extended.GetRouterOptions())

	options.FileStoreUseCase.AppendHandler(router.SubRouter("/filestore"))

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: router.CreateHttpHandler(),
	}

	return &Server{httpserver: server}
}

func (s *Server) ListenAndServe() error {
	log.Printf("Starting server on \"%s\"", s.httpserver.Addr)
	return s.httpserver.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpserver.Shutdown(ctx)
}
