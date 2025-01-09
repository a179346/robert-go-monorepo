package post_board_server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	post_board_config "github.com/a179346/robert-go-monorepo/internal/post_board/config"
	auth_use_case "github.com/a179346/robert-go-monorepo/internal/post_board/use_caes/auth"
	"github.com/a179346/robert-go-monorepo/pkg/gohf"
	"github.com/a179346/robert-go-monorepo/pkg/gohf_extended"
	"github.com/rs/cors"
)

type Options struct {
	AuthUseCase auth_use_case.AuthUseCase
}

type Server struct {
	httpserver *http.Server
}

func New(config post_board_config.ServerConfig, options Options) *Server {
	router := gohf.New()

	options.AuthUseCase.AppendHandler(router.SubRouter("/auth"))

	router.Use(gohf_extended.NotFoundHandler)

	mux := router.CreateServeMux()
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
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
