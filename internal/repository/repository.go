package repository

import (
	"github.com/rs/zerolog"
	"github.com/shanto-323/backend-scaffold/config"
	"github.com/shanto-323/backend-scaffold/internal/repository/cache"
	"github.com/shanto-323/backend-scaffold/internal/repository/database"
	"github.com/shanto-323/backend-scaffold/internal/repository/database/postgres"
	"go.opentelemetry.io/otel/trace"
)

type Repository struct {
	config *config.Config
	logger *zerolog.Logger
	tracer trace.Tracer

	DatabaseDriver database.Driver
	CacheProvider  cache.Provider
}

func New(config *config.Config, logger *zerolog.Logger, tracer trace.Tracer) (*Repository, error) {

	db, err := postgres.New(config, logger, tracer)
	if err != nil {
		return nil, err
	}

	cache, err := cache.New(config, logger, tracer)
	if err != nil {
		_ = db.Close()
		return nil, err
	}

	return &Repository{
		config: config,
		logger: logger,
		tracer: tracer,

		DatabaseDriver: db,
		CacheProvider:  cache,
	}, nil
}

func (r *Repository) Close() error {
	if err := r.DatabaseDriver.Close(); err != nil {
		return err
	}
	return nil
}
