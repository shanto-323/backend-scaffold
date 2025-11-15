package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/shanto-323/backend-scaffold/config"
	"github.com/shanto-323/backend-scaffold/internal/repository"
	"github.com/shanto-323/backend-scaffold/internal/service"
	"github.com/shanto-323/backend-scaffold/pkg/otel"
)

type Server struct {
	Config      *config.Config
	Logger      *zerolog.Logger
	Repository  *repository.Repository
	Services    *service.Services
	OTELService *otel.OtelService
	httpServer  *http.Server
}

func NewServer(logger *zerolog.Logger, config *config.Config) (*Server, error) {
	otelService, err := otel.CreateOtelService(context.Background(), config)
	if err != nil {
		return nil, err
	}

	repository, err := repository.New(config, logger, otelService)
	if err != nil {
		return nil, err
	}

	services := service.New()

	return &Server{
		Config:      config,
		Logger:      logger,
		Repository:  repository,
		Services:    services,
		OTELService: otelService,
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

func (s *Server) Stop() error {
	return nil
}
