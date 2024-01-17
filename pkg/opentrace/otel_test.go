package opentrace

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.opentelemetry.io/otel/trace"
)

type tracerKeyType int

var tracerKey tracerKeyType = 0

func WithContext(ctx context.Context, tr trace.Tracer) context.Context {
	return context.WithValue(ctx, tracerKey, tr)
}

func FronContext(ctx context.Context) trace.Tracer {
	return ctx.Value(tracerKey).(trace.Tracer)
}

func TestOTEL(t *testing.T) {
	c := &Config{
		ServiceName: "echo_100",
	}
	c.SetDefaults()
	tp := c.SetProvider()

	tr := tp.Tracer("test")

	ctx := context.Background()
	spanctx, span := tr.Start(ctx, "test")
	spanctx = WithContext(spanctx, tr)

	defer span.End()

	output("main", span)
	span.SetStatus(200, "OK")
	span.RecordError(fmt.Errorf("error"))
	SpanSetStringAttrs(span, map[string]string{
		"app": "test",
		"env": "local",
	})
	// sub1(ctx)

	// time.Sleep(5 * time.Second)
	for {
		time.Sleep(2 * time.Second)
		sub1(spanctx)
	}
}

func sub1(ctx context.Context) {

	tr := FronContext(ctx)
	_, span := tr.Start(ctx, "sub1")
	defer span.End()
	span.SetStatus(400, "Bad Request")
	output("sub1", span)
	SpanSetStringAttrs(span, map[string]string{
		"feat": "sub1",
		"env":  "local",
	})
}

func output(name string, span trace.Span) {
	tid := span.SpanContext().TraceID().String()
	sid := span.SpanContext().SpanID().String()
	fmt.Println(name, tid, sid)
}
