package delay_app_server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	delay_app_config "github.com/a179346/robert-go-monorepo/apps/delay_app/config"
	delay_use_case "github.com/a179346/robert-go-monorepo/apps/delay_app/use_caes/delay"
	"github.com/a179346/robert-go-monorepo/packages/robert_router_options"
	"github.com/a179346/robert-go-monorepo/packages/roberthttp"
)

type Options struct {
	DelayUseCase delay_use_case.DelayUseCase
}

type Server struct {
	httpserver *http.Server
}

type JsonResponse struct {
	StatusCode int         `json:"statusCode"`
	Data       interface{} `json:"data"`
}

func New(config delay_app_config.ServerConfig, options Options) *Server {
	router := roberthttp.New(robert_router_options.GetRouterOptions())

	options.DelayUseCase.HandleGroup(router.Group("/delay"))

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
