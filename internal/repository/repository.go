package repository

import (
	"github.com/rs/zerolog"
	"github.com/shanto-323/backend-scaffold/config"
	"github.com/shanto-323/backend-scaffold/internal/repository/cache"
	"github.com/shanto-323/backend-scaffold/internal/repository/database"
	"github.com/shanto-323/backend-scaffold/internal/repository/database/postgres"
	"github.com/shanto-323/backend-scaffold/pkg/otel"
)

type Repository struct {
	config      *config.Config
	logger      *zerolog.Logger
	otelService *otel.OtelService

	DatabaseDriver database.Driver
	CacheProvider  cache.Provider
}

func New(config *config.Config, logger *zerolog.Logger, otelService *otel.OtelService) (*Repository, error) {

	db, err := postgres.New(config, logger, otelService)
	if err != nil {
		return nil, err
	}

	cache, err := cache.New(config, logger, otelService)
	if err != nil {
		_ = db.Close()
		return nil, err
	}

	return &Repository{
		config:         config,
		DatabaseDriver: db,
		CacheProvider:  cache,
		logger:         logger,
		otelService:    otelService,
	}, nil
}

func (r *Repository) Close() error {
	if err := r.DatabaseDriver.Close(); err != nil {
		return err
	}
	return nil
}
