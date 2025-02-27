package delay_app_server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/a179346/robert-go-monorepo/pkg/apilog"
	"github.com/a179346/robert-go-monorepo/pkg/console"
	"github.com/a179346/robert-go-monorepo/pkg/gohf_extended"
	delay_app_config "github.com/a179346/robert-go-monorepo/services/delay_app/config"
	delay_use_case "github.com/a179346/robert-go-monorepo/services/delay_app/use_cases/delay"
	"github.com/gohf-http/gohf/v6"
	"github.com/rs/cors"
)

type Options struct {
	DelayUseCase delay_use_case.DelayUseCase
	ApiLogger    apilog.Logger
}

type Server struct {
	httpserver *http.Server
}

func New(options Options) *Server {
	router := gohf.New()

	if options.ApiLogger != nil {
		appConfig := delay_app_config.GetAppConfig()
		router.Use(gohf_extended.ApiLogMiddleware(appConfig.ID, appConfig.Version, options.ApiLogger))
	}

	router.Use(gohf_extended.RecoverMiddleware)
	router.Use(gohf.MaxBytesMiddleware(5 * 1024 * 1024))
	router.Use(gohf_extended.RequestIdMiddleware)
	router.Use(gohf_extended.ReadBodyMiddleware)

	router.GET("/healthz", gohf_extended.HealthzHandler)

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
	console.Infof("Starting server on \"%s\"", s.httpserver.Addr)
	return s.httpserver.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpserver.Shutdown(ctx)
}
