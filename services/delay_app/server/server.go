package delay_app_server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/a179346/robert-go-monorepo/pkg/apilog"
	"github.com/a179346/robert-go-monorepo/pkg/console"
	"github.com/a179346/robert-go-monorepo/pkg/gin_extended"
	delay_app_config "github.com/a179346/robert-go-monorepo/services/delay_app/config"
	delay_use_case "github.com/a179346/robert-go-monorepo/services/delay_app/use_cases/delay"
	"github.com/gin-gonic/gin"
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
	router := gin.Default()
	router.Use(gin_extended.ResponseMiddleware)

	if options.ApiLogger != nil {
		appConfig := delay_app_config.GetAppConfig()
		router.Use(gin_extended.ApiLogMiddleware(appConfig.ID, appConfig.Version, options.ApiLogger))
	}

	router.Use(gin_extended.RecoverMiddleware)
	router.Use(gin_extended.MaxBytesMiddleware(5 * 1024 * 1024))
	router.Use(gin_extended.RequestIdMiddleware)
	router.Use(gin_extended.ReadBodyMiddleware)

	router.GET("/healthz", gin_extended.HealthzHandler)

	options.DelayUseCase.AppendHandler(router.Group("/delay"))

	router.Use(gin_extended.NotFoundHandler)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", delay_app_config.GetServerConfig().Port),
		Handler: cors.AllowAll().Handler(router),
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
