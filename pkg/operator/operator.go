package operator

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

type Operator interface {
	// HandlerFunc
	Handle(ctx context.Context, arc *app.RequestContext)
}

type MethodOperator interface {
	Method() string
}

type RouteOperator interface {
	Route() string
}
