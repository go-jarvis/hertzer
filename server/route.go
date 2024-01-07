package server

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
)

type RouterGroup struct {
	path string
	r    *route.RouterGroup

	subgroups   []*RouterGroup
	middlewares []HandlerFunc
	operators   []Operator
}

func NewRouterGroup(path string) *RouterGroup {
	return &RouterGroup{
		path: path,
	}
}

func (r *RouterGroup) initialize() {
	if r.r == nil {
		return
	}

	for _, m := range r.middlewares {
		r.r.Use(m)
	}

	r.handle(r.operators...)

	for _, g := range r.subgroups {
		g.initialize()
	}
}

// Use register Middlewares
func (r *RouterGroup) Use(middleware ...HandlerFunc) {
	if len(r.middlewares) == 0 {
		r.middlewares = make([]HandlerFunc, 0)
	}

	r.middlewares = append(r.middlewares, middleware...)
}

// AppendGroup register subgroups
func (r *RouterGroup) AppendGroup(group ...*RouterGroup) {
	if len(r.subgroups) == 0 {
		r.subgroups = make([]*RouterGroup, 0)
	}
	r.subgroups = append(r.subgroups, group...)
}

// Handle register Operators
func (r *RouterGroup) Handle(opers ...Operator) {
	if len(r.operators) == 0 {
		r.operators = make([]Operator, 0)
	}
	r.operators = append(r.operators, opers...)
}

func (r *RouterGroup) handle(opers ...Operator) {
	for _, oper := range opers {
		fn := func(ctx context.Context, arc *app.RequestContext) {
			ret, err := oper.Handle(ctx, arc)
			if err != nil {
				arc.JSON(consts.StatusInternalServerError, err.Error())
				return
			}
			arc.JSON(consts.StatusOK, ret)
		}

		r.r.Handle(oper.Method(), oper.Route(), fn)
	}
}
