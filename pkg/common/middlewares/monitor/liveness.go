package monitor

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-jarvis/hertzer/pkg/httpx"
)

type Liveness struct {
	httpx.MethodGet `route:"/liveness"`
}

func (Liveness) Handle(ctx context.Context, arc *app.RequestContext) {
	arc.String(200, "ok")
}
