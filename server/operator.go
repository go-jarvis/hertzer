package server

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

type HandlerFunc = app.HandlerFunc
type Operator interface {
	Handle(ctx context.Context, arc *app.RequestContext) (any, error)
	Route() string
	Method() string
}
