package opentrace

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

type tracerType int

var currentTracerKey tracerType = 0

func TracerWithContext(ctx context.Context, tracer trace.Tracer) context.Context {
	return context.WithValue(ctx, currentTracerKey, tracer)
}

func TracerFromContext(ctx context.Context) trace.Tracer {
	return ctx.Value(currentTracerKey).(trace.Tracer)
}
