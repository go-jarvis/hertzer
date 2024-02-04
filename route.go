package hertzer

import (
	"context"
	"reflect"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/go-jarvis/hertzer/pkg/httpx"
	"github.com/go-jarvis/hertzer/pkg/operator"
	"github.com/go-jarvis/hertzer/pkg/reflectx"
)

type RouterGroup struct {
	path string
	r    *route.RouterGroup

	parent *RouterGroup // parent RouterGroup

	subgroups   []*RouterGroup
	middlewares []app.HandlerFunc
	operators   []operator.Operator
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
func (r *RouterGroup) Use(middleware ...app.HandlerFunc) {
	if len(r.middlewares) == 0 {
		r.middlewares = make([]app.HandlerFunc, 0)
	}

	r.middlewares = append(r.middlewares, middleware...)
}

// AddGroup register subgroups
func (r *RouterGroup) AddGroup(groups ...*RouterGroup) {
	if len(r.subgroups) == 0 {
		r.subgroups = make([]*RouterGroup, 0)
	}

	r.subgroups = append(r.subgroups, groups...)
}

// Handle register Operators
func (r *RouterGroup) Handle(opers ...operator.Operator) {
	if len(r.operators) == 0 {
		r.operators = make([]operator.Operator, 0)
	}
	r.operators = append(r.operators, opers...)
}

// handleFunc register app.HandlerFunc
// 1. set default content-type
// 2. bind and validate operator
// 3. change the log caller depth
func (r *RouterGroup) handlerFunc(oper operator.Operator) app.HandlerFunc {

	return func(ctx context.Context, arc *app.RequestContext) {
		// set default content-type
		v := arc.Request.Header.Get("Content-Type")
		if v == "" {
			arc.Request.Header.Set("Content-Type", "application/json")
		}

		// // Todo: caller
		// // 1. 处理 log 深度
		// log := slogr.FromContext(c)
		// if l, ok := log.(*slogr.SLogger); ok {
		// 	// l.SetCaller(l.GetCaller() + -1)
		// 	// l.SetCaller(3)
		// 	// fmt.Println("l.GetCaller()", l.GetCaller())
		// 	c = slogr.WithContext(c, l)
		// }

		// 2. 处理 handler
		err := arc.BindAndValidate(oper)
		if err != nil {
			arc.JSON(consts.StatusBadRequest, err.Error())
			return
		}

		oper.Handle(ctx, arc)
	}
}

func (r *RouterGroup) handle(opers ...operator.Operator) {
	for _, oper := range opers {
		// create a deepcopy
		oper := operator.DeepCopy(oper)

		// get method and path
		m, p := getHttpBasic(oper)

		// main: registery handler func
		hfn := r.handlerFunc(oper)

		r.r.Handle(m, p, hfn)
	}
}

func getHttpBasic(oper any) (method, path string) {

	if oper, ok := oper.(httpx.Methoder); ok {
		method = oper.Method()
	}

	if oper, ok := oper.(operator.RouteOperator); ok {
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
