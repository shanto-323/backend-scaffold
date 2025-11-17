package handler

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/shanto-323/backend-scaffold/internal/server/middleware"
	"github.com/shanto-323/backend-scaffold/internal/server/validation"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type HandlerFunc[Req validation.Validatable, Res any] func(c echo.Context, req Req) (Res, error)

type HandlerFuncNoContent[Req validation.Validatable] func(c echo.Context, req Req) error

type ResponseHandler interface {
	Handle(c echo.Context, result any) error
	GetOperation() string
}

// -------------------------- //

type JSONResponseHandler struct {
	status int
}

func (h JSONResponseHandler) Handle(c echo.Context, result any) error {
	return c.JSON(h.status, result)
}

func (h JSONResponseHandler) GetOperation() string {
	return "handler"
}

// -------------------------- //

type NoContentResponseHandler struct {
	status int
}

func (h NoContentResponseHandler) Handle(c echo.Context, result any) error {
	return c.NoContent(h.status)
}

func (h NoContentResponseHandler) GetOperation() string {
	return "handler_no_content"
}

// -------------------------- //

type FileResponseHandler struct {
	status      int
	filename    string
	contentType string
}

func (h FileResponseHandler) Handle(c echo.Context, result any) error {
	data := result.([]byte)
	c.Response().Header().Set("Content-Disposition", "attachment; filename="+h.filename)
	return c.Blob(h.status, h.contentType, data)
}

func (h FileResponseHandler) GetOperation() string {
	return "handler_file"
}

// -------------------------- //

func handleRequest[Req validation.Validatable](
	c echo.Context,
	req Req,
	handler func(c echo.Context, req Req) (any, error),
	responseHandler ResponseHandler,
) error {
	start := time.Now()
	method := c.Request().Method
	path := c.Path()
	route := path

	span := trace.SpanFromContext(c.Request().Context())

	// Get context-enhanced logger
	loggerBuilder := middleware.GetLogger(c).With().
		Str("operation", responseHandler.GetOperation()).
		Str("method", method).
		Str("path", path).
		Str("route", route)



	// Add file-specific fields to logger if it's a file handler
	if fileHandler, ok := responseHandler.(FileResponseHandler); ok {
		loggerBuilder = loggerBuilder.
			Str("filename", fileHandler.filename).
			Str("content_type", fileHandler.contentType)
	}

	logger := loggerBuilder.Logger()

	// Validation with observability
	validationStart := time.Now()
	if err := validation.BindAndValidate(c, req); err != nil {
		validationDuration := time.Since(validationStart)

		logger.Error().
			Err(err).
			Dur("validation_duration", validationDuration).
			Msg("request validation failed")
		span.RecordError(err)
		span.SetAttributes(
			attribute.String("validation.status", "failed"),
			attribute.Int64("validation.duration_ms", validationDuration.Milliseconds()),
		)
		return err
	}

	validationDuration := time.Since(validationStart)
	span.SetAttributes(
		attribute.String("validation.status", "success"),
		attribute.Int64("validation.duration_ms", validationDuration.Milliseconds()),
	)

	logger.Debug().
		Dur("validation_duration", validationDuration).
		Msg("request validation successful")

	handlerStart := time.Now()
	result, err := handler(c, req)
	handlerDuration := time.Since(handlerStart)

	if err != nil {
		totalDuration := time.Since(start)

		logger.Error().
			Err(err).
			Dur("handler_duration", handlerDuration).
			Dur("total_duration", totalDuration).
			Msg("handler execution failed")
		span.RecordError(err)
		span.SetAttributes(
			attribute.String("handler.status", "failed"),
			attribute.Int64("handler.duration_ms", handlerDuration.Milliseconds()),
			attribute.Int64("total.duration_ms", validationDuration.Milliseconds()),
		)

		return err
	}

	totalDuration := time.Since(start)
	span.SetAttributes(
		attribute.String("handler.status", "success"),
		attribute.Int64("handler.duration_ms", handlerDuration.Milliseconds()),
		attribute.Int64("total.duration_ms", totalDuration.Milliseconds()),
	)

	logger.Info().
		Dur("handler_duration", handlerDuration).
		Dur("validation_duration", validationDuration).
		Dur("total_duration", totalDuration).
		Msg("request completed successfully")

	return responseHandler.Handle(c, result)
}

// -------------------------- //

func Handle[Req validation.Validatable, Res any](
	handler HandlerFunc[Req, Res],
	status int,
	req Req,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		return handleRequest(c, req, func(c echo.Context, req Req) (any, error) {
			return handler(c, req)
		}, JSONResponseHandler{status: status})
	}
}

// -------------------------- //

func HandleFile[Req validation.Validatable](
	handler HandlerFunc[Req, []byte],
	status int,
	req Req,
	filename string,
	contentType string,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		return handleRequest(c, req, func(c echo.Context, req Req) (any, error) {
			return handler(c, req)
		}, FileResponseHandler{
			status:      status,
			filename:    filename,
			contentType: contentType,
		})
	}
}

// -------------------------- //

func HandleNoContent[Req validation.Validatable](
	handler HandlerFuncNoContent[Req],
	status int,
	req Req,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		return handleRequest(c, req, func(c echo.Context, req Req) (any, error) {
			err := handler(c, req)
			return nil, err
		}, NoContentResponseHandler{status: status})
	}
}
