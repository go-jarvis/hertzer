package logr

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-jarvis/slogr"
)

// WithContext inject slogr into context
func WithContext(ctx context.Context, log slogr.Logger) context.Context {
	return slogr.WithContext(ctx, log)
}

// FromContext exact slogr from context
func FromContext(ctx context.Context) slogr.Logger {

	log := slogr.FromContext(ctx)

	// 如果为Discard， 则返回默认的logger
	if _, ok := log.(*slogr.Discard); ok {
		return slogr.Default()
	}

	return log
}

// WithSLogger inject slogr.Logger in middleware
func WithSLogger(log slogr.Logger) app.HandlerFunc {

	return func(c context.Context, ctx *app.RequestContext) {
		c = WithContext(c, log)
		ctx.Next(c)
	}
}
