package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/shanto-323/backend-scaffold/internal/server/handler"
	"github.com/shanto-323/backend-scaffold/internal/server/middleware"
)

func RegisterV1Routes(r *echo.Group, h *handler.Handlers, m *middleware.Middlewares) {
	student := r.Group("/student")

	student.POST("", h.StudentHandler.Create)
}
