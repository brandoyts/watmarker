package server

import (
	"context"
	"net/http"

	"github.com/brandoyts/watmarker/microservice/api_gateway/config"
	"github.com/brandoyts/watmarker/pkg/common/ports/v1"
)

type Server struct {
	logger      ports.Logger
	router      *http.ServeMux
	server      *http.Server
	middlewares []Middleware
}

func NewServer(config config.GatewayConfig, logger ports.Logger) *Server {
	router := http.NewServeMux()

	httpServer := &http.Server{
		Addr:    config.Address,
		Handler: router,
	}

	return &Server{
		logger: logger,
		router: router,
		server: httpServer,
	}
}

func (s *Server) Use(m Middleware) {
	s.middlewares = append(s.middlewares, m)
}

func (s *Server) RegisterHandler(path string, handler http.HandlerFunc) {
	var wrappedHandler http.Handler = handler

	// apply middlewares in reverse order (so the first Use() wraps outermost)
	for i := len(s.middlewares) - 1; i >= 0; i-- {
		wrappedHandler = s.middlewares[i](wrappedHandler)
	}

	s.router.Handle(path, wrappedHandler)
}

func (s *Server) Run() error {
	err := s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
