package logr

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type AcceccLoggerConfig struct {
	SkipPaths []string
}

var defaultLogConfig = AcceccLoggerConfig{
	SkipPaths: []string{
		"/liveness",
		"/readiness",
		"/metrics",
		"/healthz",
		"/readyz",
	},
}

func AccessLogger() app.HandlerFunc {
	return AccessLoggerWithConfig(defaultLogConfig)
}

func AccessLoggerWithConfig(config AcceccLoggerConfig) app.HandlerFunc {
	var skip map[string]struct{}
	if length := len(config.SkipPaths); length != 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range config.SkipPaths {
			skip[path] = struct{}{}
		}
	}

	return func(ctx context.Context, ac *app.RequestContext) {

		full_path := string(ac.Request.URI().PathOriginal())
		if _, isSkip := skip[full_path]; isSkip {
			ac.Next(ctx)
			return
		}

		log := FromContext(ctx)
		ctx, log = log.WithContext(ctx, "trace_id", traceid(ctx))

		// start log
		start := time.Now()
		ac.Next(ctx) // do next handler

		end := time.Now()
		latency := end.Sub(start).Milliseconds // cost time

		log.With(
			"status", ac.Response.StatusCode(),
			"cost", fmt.Sprintf("%dms", latency()),
			"request_method", string(ac.Request.Header.Method()),
			"request_uri", full_path,
			"remote_addr", ac.ClientIP(),
			"host", string(ac.Request.Host()),
		).Info("access_log")
	}
}

// traceid return trace id from opentelemetry Trace ID
// or generate a new one by uuid
func traceid(ctx context.Context) string {
	spanctx := trace.SpanContextFromContext(ctx)

	if spanctx.IsValid() {
		if spanctx.TraceID().IsValid() {
			return spanctx.TraceID().String()
		}
	}

	return uuid.New().String()
}
