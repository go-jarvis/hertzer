package hertztrace

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/go-jarvis/hertzer/pkg/opentrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/semconv/v1.13.0/httpconv"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var defaultTracerName string = "github.com/go-jarvis/hertzer"

type config struct {
	ServiceName string
	provider    oteltrace.TracerProvider
	propagator  propagation.TextMapPropagator
}

func init() {
	otel.SetTextMapPropagator(propagation.TraceContext{})
}

func (c *config) defaults() {
	if c.ServiceName == "" {
		c.ServiceName = defaultTracerName
	}

	if c.provider == nil {
		c.provider = otel.GetTracerProvider()
	}

	if c.propagator == nil {
		c.propagator = otel.GetTextMapPropagator()
		// if c.propagator == nil {
		// 	c.propagator = propagation.TraceContext{}
		// }
	}
}

type ServerOption func(*config)

func WithServiceName(name string) ServerOption {
	return func(c *config) {
		c.ServiceName = name
	}
}

func WithPropagator(p propagation.TextMapPropagator) ServerOption {
	return func(c *config) {
		c.propagator = p
	}
}

func ServerMiddleware(opts ...ServerOption) app.HandlerFunc {
	cfg := &config{}
	cfg.defaults()

	for _, opt := range opts {
		opt(cfg)
	}

	tracer := otel.GetTracerProvider().Tracer(cfg.ServiceName)

	return func(ctx context.Context, ac *app.RequestContext) {
		ctx = opentrace.TracerWithContext(ctx, tracer)

		spanOpts := []oteltrace.SpanStartOption{
			oteltrace.WithSpanKind(oteltrace.SpanKindServer),
			oteltrace.WithAttributes(basicAttrs(ac)...),
		}

		// extract baggage and span context from header
		bags, spanCtx := Extract(ctx, cfg, &ac.Request.Header)

		// set baggage
		ctx = baggage.ContextWithBaggage(ctx, bags)

		ctx = oteltrace.ContextWithRemoteSpanContext(ctx, spanCtx)
		spanName := ac.FullPath()

		spanctx, span := opentrace.SpanStart(ctx, spanName, spanOpts...)
		defer span.End()

		ac.Next(spanctx)

		code := ac.Response.StatusCode()
		if code/100 == 2 {
			span.SetStatus(codes.Ok, consts.StatusMessage(code))
		} else {
			span.SetStatus(httpconv.ServerStatus(code))
		}

		responseTraceparentInject(spanctx, ac)
	}
}

// basicAttrs 添加 http 请求基础属性
func basicAttrs(ac *app.RequestContext) []attribute.KeyValue {
	attrs := []attribute.KeyValue{}
	attrs = append(attrs, attribute.String("http.method", string(ac.Method())))
	attrs = append(attrs, attribute.String("http.path", ac.FullPath()))
	attrs = append(attrs, attribute.String("http.host", string(ac.Host())))
	attrs = append(attrs, attribute.String("http.scheme", string(ac.GetRequest().Scheme())))
	return attrs
}

// responseTraceparentInject 向后传递 Header: traceparent
func responseTraceparentInject(ctx context.Context, ac *app.RequestContext) {
	// 6. 向后传递 Header: traceparent
	pp := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
	)

	carrier := propagation.MapCarrier{}
	pp.Inject(ctx, carrier)

	for k, v := range carrier {
		ac.Response.Header.Set(k, v)
	}
}
