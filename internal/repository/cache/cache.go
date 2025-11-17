package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/shanto-323/backend-scaffold/config"
	"go.opentelemetry.io/otel/trace"
)

type Provider interface {
	Close() error
	Ping(ctx context.Context) error
}

type cache struct {
	logger *zerolog.Logger
	Client *redis.Client
}

func New(config *config.Config, logger *zerolog.Logger, tracer trace.Tracer) (Provider, error) {
	if config == nil || logger == nil {
		return nil, fmt.Errorf("config and logger must not be nil")
	}

	opt, _ := redis.ParseURL(config.Redis.Address)

	redisClient := redis.NewClient(opt)

	// Validate connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		_ = redisClient.Close()
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	logger.Info().Msg("redis service initialized successfully")

	return &cache{
		logger: logger,
		Client: redisClient,
	}, nil
}

func (c *cache) Ping(ctx context.Context) error {
	return c.Client.Ping(ctx).Err()
}

// Close gracefully closes the Redis connection
func (c *cache) Close() error {
	if err := c.Client.Close(); err != nil {
		c.logger.Error().Err(err).Msg("Error closing Redis connection")
		return err
	}

	c.logger.Info().Msg("Redis connection closed")
	return nil
}
