package opentrace

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// SpanStart starts a trace with the given name and options.
// and returns the span context and span.
func SpanStart(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {

	tr := TracerFromContext(ctx)
	spanctx, span := tr.Start(ctx, name, opts...)

	return spanctx, span
}

// SpanEnd ends the span.
// use defer SpanEnd(span) to end the span.
func SpanEnd(span trace.Span) {
	if span != nil {
		span.End()
	}
}

// SpanSetStringAttrs sets the attributes to the span.
func SpanSetStringAttrs(span trace.Span, kvs map[string]string) {
	attrkv := []attribute.KeyValue{}

	for k, v := range kvs {
		attrkv = append(attrkv, attribute.KeyValue{
			Key:   attribute.Key(k),
			Value: attribute.StringValue(v),
		})
	}

	span.SetAttributes(attrkv...)
}

type spanAttrType int

var currentAttrKey spanAttrType = 0

// SpanAttrWithContext sets the attributes to the span in context.
func SpanAttrWithContext(ctx context.Context, kvs map[string]string) context.Context {

	m := SpanAttrFromContext(ctx)

	for k, v := range kvs {
		m[k] = v
	}

	return context.WithValue(ctx, currentAttrKey, m)
}

// SpanAttrFromContext gets the attributes from the span in context.
func SpanAttrFromContext(ctx context.Context) map[string]string {
	m, ok := ctx.Value(currentAttrKey).(map[string]string)
	if ok {
		return m
	}
	return make(map[string]string, 0)
}
