package calculator

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/a179346/robert-go-monorepo/pkg/console"
	"github.com/a179346/robert-go-monorepo/pkg/jsonvalidator"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/rs/cors"
)

type httpServer struct {
	httpserver *http.Server
}

func NewHttpServer(port uint, endpoints Endpoints) *httpServer {
	mux := http.NewServeMux()

	mux.Handle("POST /divide", httptransport.NewServer(
		endpoints.divide,
		httpDecodeDivideRequest,
		httpEncodeResponse,
	))

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: cors.AllowAll().Handler(mux),
	}

	return &httpServer{httpserver: server}
}

func (s *httpServer) ListenAndServe() error {
	console.Infof("Starting server on \"%s\"", s.httpserver.Addr)
	return s.httpserver.ListenAndServe()
}

func (s *httpServer) Shutdown(ctx context.Context) error {
	return s.httpserver.Shutdown(ctx)
}

func httpDecodeDivideRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return jsonvalidator.FromReader[divideReq](r.Body)
}

func httpEncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
