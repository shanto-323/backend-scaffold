package router

import (
	"github.com/labstack/echo/v4"
	"github.com/shanto-323/backend-scaffold/internal/server/handler"
)

func registerSystemRouter(r *echo.Echo, h *handler.HealthHandler) {
	r.GET("/status", h.CheckHealth)
}
