package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/shanto-323/backend-scaffold/internal/server"
	"go.opentelemetry.io/otel/attribute"
)

type Tracer struct {
	s *server.Server
}

func NewTracer(s *server.Server) *Tracer {
	return &Tracer{
		s: s,
	}
}

func (t *Tracer) EnhanceTracing() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()

			ctx, span := t.s.TraceProvider.Tracer.Start(ctx, c.Path())
			c.SetRequest(c.Request().WithContext(ctx))

			attrs := []attribute.KeyValue{}
			request_id := GetRequestID(c)
			if request_id != "" {
				attrs = append(attrs, attribute.String("request_id", request_id))
			}

			user_id := GetUserID(c)
			if user_id != "" {
				attrs = append(attrs, attribute.String("user_id", user_id))
			}

			attrs = append(
				attrs,
				attribute.String("http.method", c.Request().Method),
				attribute.String("http.path", c.Path()),
				attribute.String("http.ip", c.RealIP()),
			)

			span.SetAttributes(attrs...)

			err := next(c)
			if err != nil {
				span.RecordError(err)
			}

			span.SetAttributes(attribute.Int("http.status_code", c.Response().Status))
			span.End()

			if err := t.s.TraceProvider.ForceFlush(ctx); err != nil {
				t.s.Logger.Error().Err(err).Msg("Failed to flush span")
			}

			return err
		}
	}
}
