package operator

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

type Operator interface {
	Handle(ctx context.Context, arc *app.RequestContext) (any, error)
}

type PreHandlersOperator interface {
	PreHandlers() []app.HandlerFunc
}

type PostHandlersOperator interface {
	PostHandlers() []app.HandlerFunc
}

type MethodOperator interface {
	Method() string
}

type RouteOperator interface {
	Route() string
}
