package main

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-jarvis/hertzer"
	"github.com/go-jarvis/hertzer/pkg/common/middlewares/hertztrace"
	"github.com/go-jarvis/hertzer/pkg/httpx"
	"github.com/go-jarvis/hertzer/pkg/opentrace"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	tt := &opentrace.Config{
		Endpoint:    OtelEndpoint,
		ServiceName: "hertzer-demo",
		BasicAuth:   OtelAuth,
	}
	tt.SetDefaults()
	_ = tt.SetProvider()

	s := &hertzer.Server{}
	s.Use(
		hertztrace.ServerMiddleware(),
	)
	s.Handle(&PingPong{})
	s.Run()
}

type PingPong struct {
	httpx.MethodGet `route:"/ping"`
}

func (PingPong) Handle(ctx context.Context, arc *app.RequestContext) (any, error) {

	name := "ping-pong"
	spanctx, span := opentrace.SpanStart(ctx, name)
	defer StopSpan(span)

	fmt.Println(name, span.SpanContext().TraceID().String(), span.SpanContext().SpanID().String())

	do(spanctx)
	do2(spanctx)

	return "pong", nil
}

func do(ctx context.Context) {
	name := "do"
	_, span := opentrace.SpanStart(ctx, name)
	defer StopSpan(span)

	ctx = opentrace.SpanAttrWithContext(ctx, map[string]string{
		"key1": "value1",
	})
	opentrace.SpanSetStringAttrs(span, opentrace.SpanAttrFromContext(ctx))

	fmt.Println(name, span.SpanContext().TraceID().String(), span.SpanContext().SpanID().String())
}

func StopSpan(span trace.Span) {
	if span != nil {
		span.End()
	}
}

func do2(ctx context.Context) {
	name := "do2"
	_, span := opentrace.SpanStart(ctx, name)
	defer StopSpan(span)

	ctx = opentrace.SpanAttrWithContext(ctx, map[string]string{
		"key2": "value2",
	})
	opentrace.SpanSetStringAttrs(span, opentrace.SpanAttrFromContext(ctx))

	fmt.Println(name, span.SpanContext().TraceID().String(), span.SpanContext().SpanID().String())
}
