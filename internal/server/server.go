package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/shanto-323/backend-scaffold/config"
	"github.com/shanto-323/backend-scaffold/internal/repository"
	"github.com/shanto-323/backend-scaffold/pkg/tracer"
)

type Server struct {
	Config        *config.Config
	Logger        *zerolog.Logger
	Repository    *repository.Repository
	TraceProvider *tracer.TraceProvider
	httpServer    *http.Server
}

func NewServer(logger *zerolog.Logger, config *config.Config) (*Server, error) {
	logger.Info().Msg(config.Monitor.OTEL.TempoEndpoint)
	tp, err := tracer.New(context.Background(), config)
	if err != nil {
		return nil, err
	}

	repository, err := repository.New(config, logger, tp.Tracer)
	if err != nil {
		return nil, err
	}


	return &Server{
		Config:        config,
		Logger:        logger,
		Repository:    repository,
		TraceProvider: tp,
	}, nil
}

func (s *Server) SetUpHTTPServer(handler http.Handler) {
	s.httpServer = &http.Server{
		Addr:         ":" + s.Config.Server.Port,
		Handler:      handler,
		ReadTimeout:  time.Duration(s.Config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.Config.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(s.Config.Server.IdleTimeout) * time.Second,
	}
}

func (s *Server) Run() error {
	if s.httpServer == nil {
		return errors.New("HTTP server not initialized")
	}

	s.Logger.Info().
		Str("port", s.Config.Server.Port).
		Str("env", s.Config.Primary.Env).
		Msg("starting server")
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.TraceProvider.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
