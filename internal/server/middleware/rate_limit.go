package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/shanto-323/backend-scaffold/internal/server"
)

type RateLimit struct {
	s *server.Server
}

func NewRateLimit(s *server.Server) *RateLimit {
	return &RateLimit{
		s: s,
	}
}

func (r *RateLimit) RateLimitHit() echo.MiddlewareFunc {
	return middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Store: middleware.NewRateLimiterMemoryStore(10),
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return echo.NewHTTPError(http.StatusTooManyRequests, "Rate limit exceeded")
		},
	})
}
