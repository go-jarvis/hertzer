package server

import (
	"context"
	"reflect"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/go-jarvis/herts/pkg/httpx"
	"github.com/go-jarvis/herts/pkg/reflectx"
)

type RouterGroup struct {
	path string
	r    *route.RouterGroup

	parent *RouterGroup // parent RouterGroup

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
		r.r = r.parent.r.Group(r.path)
	}

	for _, m := range r.middlewares {
		r.r.Use(m)
	}

	r.handle(r.operators...)

	for _, sub := range r.subgroups {
		// set parent as self
		sub.parent = r

		// initialize subgroups
		sub.initialize()
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
func (r *RouterGroup) AppendGroup(groups ...*RouterGroup) {
	if len(r.subgroups) == 0 {
		r.subgroups = make([]*RouterGroup, 0)
	}

	r.subgroups = append(r.subgroups, groups...)
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

			v := arc.Request.Header.Get("Content-Type")
			if v == "" {
				arc.Request.Header.Set("Content-Type", "application/json")
			}

			err := arc.Bind(oper)
			if err != nil {
				arc.JSON(consts.StatusBadRequest, err.Error())
				return
			}

			ret, err := oper.Handle(ctx, arc)
			if err != nil {
				arc.JSON(consts.StatusInternalServerError, err.Error())
				return
			}
			arc.JSON(consts.StatusOK, ret)
		}

		m, p := getHttpBasic(oper)
		r.r.Handle(m, p, fn)
	}
}

func getHttpBasic(oper any) (method, path string) {

	if oper, ok := oper.(httpx.Methoder); ok {
		method = oper.Method()
	}

	if oper, ok := oper.(Router); ok {
		path = oper.Route()
	}

	if path != "" && method != "" {
		return method, path
	}

	// get path by use reflect
	rt := reflect.TypeOf(oper)
	rt = reflectx.Deref(rt)
	for i := 0; i < rt.NumField(); i++ {
		ft := rt.Field(i)
		// 取一个
		val, ok := ft.Tag.Lookup("route")
		if ok {
			path = val
			break
		}
	}

	return
}
