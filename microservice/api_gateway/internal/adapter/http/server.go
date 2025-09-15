package http

import (
	"context"
	"net/http"

	"github.com/brandoyts/watmarker/microservice/api_gateway/internal/adapter/http/middleware"
)

type Server struct {
	router      *http.ServeMux
	server      *http.Server
	middlewares []middleware.Middleware
}

func NewServer(address string) *Server {
	router := http.NewServeMux()

	return &Server{
		router: router,
		server: &http.Server{
			Addr:    address,
			Handler: router,
		},
	}
}

func (s *Server) Use(m middleware.Middleware) {
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
