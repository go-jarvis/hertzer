package logr

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-jarvis/slogr"
)

type slogrKeytype int

var slogrKey slogrKeytype = 0

// withContext inject slogr into context
func withContext(ctx context.Context, log slogr.Logger) context.Context {
	return context.WithValue(ctx, slogrKey, log)
}

// FromContext exact slogr from context
func FromContext(ctx context.Context) slogr.Logger {
	log, ok := ctx.Value(slogrKey).(slogr.Logger)
	if ok {
		return log
	}
	return slogr.Default()
}

// WithSLogger inject slogr.Logger in middleware
func WithSLogger(log slogr.Logger) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		c = withContext(c, log)
		ctx.Next(c)
	}
}
