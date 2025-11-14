package otel

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/shanto-323/backend-scaffold/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"
)

type OtelService struct {
	Tracer         trace.Tracer
	tracerProvider *tracesdk.TracerProvider
}

func CreateOtelService(ctx context.Context, config *config.Config) (*OtelService, error) {
	otelService := &OtelService{}

	// Create resource
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(config.Monitor.ServiceName),
			semconv.DeploymentEnvironment(config.Primary.Env),
		),
	)
	if err != nil {
		return nil, err
	}

	// Setup trace exporter to Tempo
	traceExporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(config.Monitor.OTEL.TempoEndpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	// Create tracer provider
	tracerProvider := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(traceExporter),
		tracesdk.WithResource(res),
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
	)

	otel.SetTracerProvider(tracerProvider)
	otelService.tracerProvider = tracerProvider
	otelService.Tracer = tracerProvider.Tracer(config.Monitor.ServiceName)

	return otelService, nil
}

func (os *OtelService) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		// Start trace span
		ctx, span := os.Tracer.Start(ctx, c.Request().Method+" "+c.Request().URL.Path,
			trace.WithAttributes(
				attribute.String("http.method", c.Request().Method),
				attribute.String("http.url", c.Request().URL.String()),
				attribute.String("http.scheme", c.Request().URL.Scheme),
				attribute.String("http.target", c.Request().URL.Path),
				attribute.String("http.host", c.Request().Host),
			),
		)
		defer span.End()

		// Execute handler
		err := next(c)

		// Set response status
		span.SetAttributes(
			attribute.Int("http.status_code", c.Response().Status),
		)

		if err != nil {
			span.RecordError(err)
		}

		return err
	}
}

func (os *OtelService) Shutdown(ctx context.Context) error {
	if err := os.tracerProvider.ForceFlush(ctx); err != nil {
		return err
	}
	return os.tracerProvider.Shutdown(ctx)
}

