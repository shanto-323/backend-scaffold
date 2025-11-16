package tracer

import (
	"context"

	"github.com/shanto-323/backend-scaffold/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace"
)

type TraceProvider struct {
	traceProvider *tracesdk.TracerProvider
	Tracer        trace.Tracer
}

func New(ctx context.Context, config *config.Config) (*TraceProvider, error) {
	exp, err := otlptracehttp.New(ctx,
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithEndpoint(config.Monitor.OTEL.TempoEndpoint),
	)
	if err != nil {
		return nil, err
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(config.Monitor.ServiceName),
			semconv.DeploymentEnvironmentName(config.Primary.Env),
		),
	)
	if err != nil {
		return nil, err
	}

	consoleExp, _ := stdouttrace.New()

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithSyncer(consoleExp),
		tracesdk.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	tracer := tp.Tracer(config.Monitor.ServiceName)
	return &TraceProvider{
		traceProvider: tp,
		Tracer:        tracer,
	}, nil
}

func (tp *TraceProvider) ForceFlush(ctx context.Context) error {
	return tp.traceProvider.ForceFlush(ctx)
}

func (tp *TraceProvider) Shutdown(ctx context.Context) error {
	return tp.traceProvider.Shutdown(ctx)
}
