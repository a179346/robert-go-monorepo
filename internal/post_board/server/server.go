package post_board_server

import (
	"context"
	"fmt"
	"net/http"

	post_board_config "github.com/a179346/robert-go-monorepo/internal/post_board/config"
	"github.com/a179346/robert-go-monorepo/internal/post_board/middlewares"
	auth_use_case "github.com/a179346/robert-go-monorepo/internal/post_board/use_cases/auth"
	post_use_case "github.com/a179346/robert-go-monorepo/internal/post_board/use_cases/post"
	user_use_case "github.com/a179346/robert-go-monorepo/internal/post_board/use_cases/user"
	"github.com/a179346/robert-go-monorepo/pkg/console"
	"github.com/a179346/robert-go-monorepo/pkg/gohf_extended"
	"github.com/gohf-http/gohf/v6"
	"github.com/rs/cors"
)

type Options struct {
	AuthUseCase auth_use_case.AuthUseCase
	UserUseCase user_use_case.UserUseCase
	PostUseCase post_use_case.PostUseCase
	ApiLogger   gohf_extended.ApiLogger
}

type Server struct {
	httpserver *http.Server
}

func New(options Options) *Server {
	router := gohf.New()

	if options.ApiLogger != nil {
		appConfig := post_board_config.GetAppConfig()
		router.Use(gohf_extended.ApiLogMiddleware(appConfig.ID, appConfig.Version, options.ApiLogger))
	}

	router.Use(gohf.MaxBytesMiddleware(5 * 1024 * 1024))
	router.Use(gohf_extended.RequestIdMiddleware)
	router.Use(gohf_extended.ReadBodyMiddleware)

	router.GET("/healthz", gohf_extended.HealthzHandler)

	options.AuthUseCase.AppendHandler(router.SubRouter("/auth"))

	{
		authedRouter := router.SubRouter("/authed")
		authedRouter.Use(middlewares.AuthedMiddleware)

		options.UserUseCase.AppendHandler(authedRouter.SubRouter("/users"))
		options.PostUseCase.AppendHandler(authedRouter.SubRouter("/posts"))
	}

	router.Use(gohf_extended.NotFoundHandler)

	mux := router.CreateServeMux()
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", post_board_config.GetServerConfig().Port),
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
