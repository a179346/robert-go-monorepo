package delay_app_server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	delay_app_config "github.com/a179346/robert-go-monorepo/apps/delay_app/config"
	delay_use_case "github.com/a179346/robert-go-monorepo/apps/delay_app/use_caes/delay"
)

type Options struct {
	DelayUseCase delay_use_case.DelayUseCase
}

type Server struct {
	httpserver *http.Server
}

func New(config delay_app_config.ServerConfig, options Options) *Server {
	mux := http.NewServeMux()

	options.DelayUseCase.AddRoutesTo(mux)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: mux,
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