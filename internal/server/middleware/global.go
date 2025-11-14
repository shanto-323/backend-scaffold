package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/shanto-323/backend-scaffold/internal/server"
)

type Global struct {
	s *server.Server
}

func NewGlobal(s *server.Server) *Global {
	return &Global{
		s: s,
	}
}

func (g *Global) CROS() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: g.s.Config.Server.CORSAllowedOrigins,
	})
}

