package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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
	go func() {
		s.logger.Info("Starting API Gateway server on ", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("Server failed to start: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	s.logger.Info("Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), s.server.IdleTimeout)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}

	s.logger.Info("Server exited successfully.")
	return nil
}
