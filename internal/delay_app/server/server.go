package delay_app_server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	delay_app_config "github.com/a179346/robert-go-monorepo/internal/delay_app/config"
	delay_use_case "github.com/a179346/robert-go-monorepo/internal/delay_app/use_cases/delay"
	"github.com/a179346/robert-go-monorepo/pkg/gohf_extended"
	"github.com/gohf-http/gohf/v2"
	"github.com/rs/cors"
)

type Options struct {
	DelayUseCase delay_use_case.DelayUseCase
}

type Server struct {
	httpserver *http.Server
}

func New(options Options) *Server {
	router := gohf.New()

	options.DelayUseCase.AppendHandler(router.SubRouter("/delay"))

	router.Use(gohf_extended.NotFoundHandler)

	mux := router.CreateServeMux()
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", delay_app_config.GetServerConfig().Port),
		Handler: cors.AllowAll().Handler(mux),
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
