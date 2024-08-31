package server

import (
	"context"
	"errors"
	"github.com/askaroe/reservationAPI/pkg/jsonlog"
	"github.com/go-chi/chi/v5"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	httpServer *http.Server
	logger     *jsonlog.Logger
}

func NewServer(router chi.Router, logger *jsonlog.Logger) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    ":8080",
			Handler: router,
		},
		logger: logger,
	}
}

func (s *Server) Start() {
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.PrintFatal(err, nil)
		}
	}()
	s.logger.PrintInfo("server started", nil)
}

func (s *Server) Shutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	s.logger.PrintInfo("shutting down server...", nil)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.PrintFatal(err, nil)
	}
	s.logger.PrintInfo("server exiting", nil)
}
