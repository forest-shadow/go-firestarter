package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/forest-shadow/go-firestarter/pkg/logger"
)

type Server struct {
	server *http.Server
	config Config
	log    logger.Logger
}

func New(handler http.Handler, c Config, log logger.Logger) (*Server, error) {
	c = c.WithDefaults()
	if err := c.Validate(); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	if handler == nil {
		return nil, errors.New("handler is required")
	}

	if log == nil {
		return nil, errors.New("logger is required")
	}

	httpServer := &http.Server{
		Handler:           handler,
		ReadTimeout:       c.ReadTimeout,
		ReadHeaderTimeout: c.ReadHeaderTimeout,
		WriteTimeout:      c.WriteTimeout,
		IdleTimeout:       c.IdleTimeout,
		Addr:              net.JoinHostPort("", c.Port),
	}

	return &Server{
		server: httpServer,
		config: c,
		log:    log,
	}, nil
}

func (s *Server) Run() error {
	listener, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		return fmt.Errorf("listen on %q: %w", s.server.Addr, err)
	}

	s.log.Info("http server: started", logger.F("address", listener.Addr().String()))

	err = s.server.Serve(listener)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("serve HTTP: %w", err)
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, s.config.ShutdownTimeout)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown HTTP server: %w", err)
	}

	s.log.Info("http server: closed")

	return nil
}
