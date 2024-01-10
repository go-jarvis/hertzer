package logr

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
)

type AcceccLogConfig struct {
	SkipPaths []string
}

func AccessLogWithConfig(config AcceccLogConfig) app.HandlerFunc {
	var skip map[string]struct{}
	if length := len(config.SkipPaths); length != 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range config.SkipPaths {
			skip[path] = struct{}{}
		}
	}

	return func(c context.Context, ctx *app.RequestContext) {

		full_path := string(ctx.Request.URI().PathOriginal())
		if _, isSkip := skip[full_path]; isSkip {
			ctx.Next(c)
			return
		}

		// start log
		start := time.Now()
		ctx.Next(c)
		end := time.Now()

		latency := end.Sub(start).Microseconds

		log := FromContext(c)
		log.With(
			"status", ctx.Response.StatusCode(),
			"cost", fmt.Sprintf("%dms", latency()),
			"method", string(ctx.Request.Header.Method()),
			"full_path", full_path,
			"client_ip", ctx.ClientIP(),
			"host", string(ctx.Request.Host()),
		).Info("access_log")
	}
}
