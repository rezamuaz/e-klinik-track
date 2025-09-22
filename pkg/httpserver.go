package pkg

// Package httpserver implements HTTP Server.

import (
	"context"
	"e-klinik/config"
	"errors"

	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	_defaultReadTimeout     = 5 * time.Second
	_defaultWriteTimeout    = 5 * time.Second
	_defaultAddr            = ":80"
	_defaultShutdownTimeout = 3 * time.Second
)

// Server -.
type Server struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
	Router          *gin.Engine
}

func NewHttp(cfg *config.Config, router *gin.Engine) *Server {

	s := prepareHttpServer(cfg, router)
	s.start()

	return s
}

func prepareHttpServer(cfg *config.Config, router *gin.Engine) *Server {
	httpServer := &http.Server{
		Handler:      router,
		ReadTimeout:  _defaultReadTimeout,
		WriteTimeout: _defaultWriteTimeout,
		Addr:         cfg.Server.ExternalPort,
	}
	httpServer.Addr = net.JoinHostPort("", cfg.Server.ExternalPort)

	s := &Server{
		server:          httpServer,
		notify:          make(chan error, 1),
		shutdownTimeout: _defaultShutdownTimeout,
		Router:          router,
	}
	return s
}

// Start
func (s *Server) start() {
	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()
}

// Notify -.
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown -.
func (s *Server) Shutdown() error {
	if s.server == nil {
		return errors.New("server is not initialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
