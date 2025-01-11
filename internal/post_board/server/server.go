package post_board_server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	post_board_config "github.com/a179346/robert-go-monorepo/internal/post_board/config"
	"github.com/a179346/robert-go-monorepo/internal/post_board/middlewares"
	auth_use_case "github.com/a179346/robert-go-monorepo/internal/post_board/use_cases/auth"
	user_use_case "github.com/a179346/robert-go-monorepo/internal/post_board/use_cases/user"
	"github.com/a179346/robert-go-monorepo/pkg/gohf_extended"
	"github.com/gohf-http/gohf/v3"
	"github.com/rs/cors"
)

type Options struct {
	AuthUseCase auth_use_case.AuthUseCase
	UserUseCase user_use_case.UserUseCase
}

type Server struct {
	httpserver *http.Server
}

func New(options Options) *Server {
	router := gohf.New()

	options.AuthUseCase.AppendHandler(router.SubRouter("/auth"))

	{
		authedRouter := router.SubRouter("/authed")
		authedRouter.Use(middlewares.AuthedMiddleware)

		options.UserUseCase.AppendHandler(authedRouter.SubRouter("/users"))
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
	log.Printf("Starting server on \"%s\"", s.httpserver.Addr)
	return s.httpserver.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpserver.Shutdown(ctx)
}
