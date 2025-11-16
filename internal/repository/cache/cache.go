package cache

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/shanto-323/backend-scaffold/config"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Provider interface {
	Close() error
}

type multiHook struct {
	hooks []redis.Hook
	mu    sync.RWMutex
}

type redisHook struct {
	logger *zerolog.Logger
	tracer trace.Tracer
}

type cache struct {
	Logger *zerolog.Logger
	Client *redis.Client
}

// DialHook for multi-hook
func (mh *multiHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		for _, hook := range mh.hooks {
			next = hook.DialHook(next)
		}
		return next(ctx, network, addr)
	}
}

// ProcessHook for multi-hook
func (mh *multiHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		for _, hook := range mh.hooks {
			next = hook.ProcessHook(next)
		}
		return next(ctx, cmd)
	}
}

// ProcessPipelineHook for multi-hook
func (mh *multiHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		for _, hook := range mh.hooks {
			next = hook.ProcessPipelineHook(next)
		}
		return next(ctx, cmds)
	}
}

// DialHook for redisHook
func (h *redisHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		h.logger.Debug().Str("network", network).Str("addr", addr).Msg("redis dial")
		return next(ctx, network, addr)
	}
}

// ProcessHook with optional tracer support
func (h *redisHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		if h.tracer != nil {
			// Development: use tracing
			tracer, ok := ctx.Value("tracer").(trace.Tracer)
			if ok {
				var span trace.Span
				ctx, span = tracer.Start(ctx, cmd.Name())
				defer span.End()

				span.SetAttributes(
					attribute.String("redis.command", cmd.Name()),
					attribute.String("redis.args", fmt.Sprint(cmd.Args())),
				)

				// Execute command
				err := next(ctx, cmd)

				// Record error if any
				if cmdErr := cmd.Err(); cmdErr != nil && cmdErr != redis.Nil {
					span.RecordError(cmdErr)
					span.SetAttributes(attribute.String("redis.error", cmdErr.Error()))
					h.logger.Error().Err(cmdErr).Str("command", cmd.Name()).Msg("redis command failed")
				}

				return err
			}
		}

		// Production: only use logger
		h.logger.Debug().Str("command", cmd.Name()).Str("args", fmt.Sprint(cmd.Args())).Msg("executing redis command")
		err := next(ctx, cmd)

		if cmdErr := cmd.Err(); cmdErr != nil && cmdErr != redis.Nil {
			h.logger.Error().Err(cmdErr).Str("command", cmd.Name()).Msg("redis command failed")
		}

		return err
	}
}

// ProcessPipelineHook
func (h *redisHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		h.logger.Debug().Int("pipeline_size", len(cmds)).Msg("executing redis pipeline")

		err := next(ctx, cmds)

		for _, cmd := range cmds {
			if cmdErr := cmd.Err(); cmdErr != nil && cmdErr != redis.Nil {
				h.logger.Error().Err(cmdErr).Str("command", cmd.Name()).Msg("redis pipeline command failed")
			}
		}

		return err
	}
}

// New creates and returns a new Redis cache provider with hooks
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
		redisClient.Close()
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	// Create hooks based on environment and otelService
	var hooks []redis.Hook

	// Always add the main redis hook with logger

	mainHook := &redisHook{
		logger: logger,
		tracer: tracer,
	}
	hooks = append(hooks, mainHook)


	// Register hooks
	if len(hooks) == 1 {
		redisClient.AddHook(hooks[0])
	} else if len(hooks) > 1 {
		multiHook := &multiHook{hooks: hooks}
		redisClient.AddHook(multiHook)
	}

	logger.Info().Msg("redis service initialized successfully")

	return &cache{
		Logger: logger,
		Client: redisClient,
	}, nil
}

// Close gracefully closes the Redis connection
func (c *cache) Close() error {
	if err := c.Client.Close(); err != nil {
		c.Logger.Error().Err(err).Msg("Error closing Redis connection")
		return err
	}

	c.Logger.Info().Msg("Redis connection closed")
	return nil
}
