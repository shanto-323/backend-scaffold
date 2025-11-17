package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/shanto-323/backend-scaffold/internal/server"
	"github.com/shanto-323/backend-scaffold/internal/server/middleware"
	"github.com/shanto-323/backend-scaffold/model"
)

const (
	Healthy   = "healthy"
	Unhealthy = "unhealthy"
)

type HealthHandler struct {
	server *server.Server
}

func NewHealthHandler(s *server.Server) *HealthHandler {
	return &HealthHandler{
		server: s,
	}
}

func (h *HealthHandler) CheckHealth(c echo.Context) error {
	start := time.Now()
	logger := middleware.GetLogger(c).With().
		Str("operation", "health_check").
		Logger()

	response := &model.Report{
		Status:      Healthy,
		Timestamp:   time.Now().UTC(),
		Environment: h.server.Config.Primary.Env,
	}

	isHealthy := true
	checks := []model.Check{}

	if h.server.Repository.DatabaseDriver != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		dbCheck := model.Check{
			Name: "postgres",
		}
		dbStart := time.Now()
		if err := h.server.Repository.DatabaseDriver.Ping(ctx); err != nil {
			isHealthy = false
			dbCheck.Status = Unhealthy
			dbCheck.ResponseTime = time.Since(dbStart).String()
			dbCheck.Error = err.Error()
			logger.Error().
				Str("check_type", "postgres").
				Str("operation", "health_check").
				Str("error_type", "postgres_unhealthy").
				Int64("response_time_ms", time.Since(dbStart).Milliseconds()).
				Str("error_message", err.Error()).
				Msg("HealthCheckError")
		} else {
			dbCheck.Status = Healthy
			dbCheck.ResponseTime = time.Since(dbStart).String()
			logger.Info().
				Dur("response_time", time.Since(dbStart)).
				Msg("database health check passed")
		}

		checks = append(checks, dbCheck)
	}

	if h.server.Repository.CacheProvider != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		redisCheck := model.Check{
			Name: "redis",
		}

		redisStart := time.Now()
		if err := h.server.Repository.CacheProvider.Ping(ctx); err != nil {
			isHealthy = false
			redisCheck.Status = Unhealthy
			redisCheck.ResponseTime = time.Since(redisStart).String()
			redisCheck.Error = err.Error()
			logger.Error().
				Str("check_type", "redis").
				Str("operation", "health_check").
				Str("error_type", "redis_unhealthy").
				Int64("response_time_ms", time.Since(redisStart).Milliseconds()).
				Str("error_message", err.Error()).
				Msg("HealthCheckError")
		} else {
			redisCheck.Status = Healthy
			redisCheck.ResponseTime = time.Since(redisStart).String()
			logger.Info().
				Dur("response_time", time.Since(redisStart)).
				Msg("redis health check passed")
		}

		checks = append(checks, redisCheck)
	}

	response.Checks = checks

	if !isHealthy {
		response.Status = Unhealthy
		logger.Error().
			Str("check_type", "overall").
			Str("operation", "health_check").
			Str("error_type", "overall_unhealthy").
			Int64("total_duration_ms", time.Since(start).Milliseconds()).
			Msg("HealthCheckError")
		return c.JSON(http.StatusServiceUnavailable, response)
	}

	logger.Info().
		Dur("total_duration", time.Since(start)).
		Msg("health check passed")

	err := c.JSON(http.StatusOK, response)
	if err != nil {
		logger.Error().
			Str("check_type", "response").
			Str("operation", "health_check").
			Str("error_type", "json_response_error").
			Str("error_message", err.Error()).
			Msg("HealthCheckError")
		return fmt.Errorf("failed to write JSON response: %w", err)
	}

	return nil
}
