package main

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/go-jarvis/hertzer"
	"github.com/go-jarvis/hertzer/pkg/common/resp"
	"github.com/go-jarvis/hertzer/pkg/httpx"
)

func main() {

	// err := errors.ErrBadPoolConn

	s := &hertzer.Server{
		Listen: ":8081",
	}

	s.Use(prefunc(), postfunc())

	s.WithOptions(
		server.WithBasePath("/api"),
		server.WithIdleTimeout(10),
	)

	s.Handle(&Ping{})

	v1 := hertzer.NewRouterGroup("/v1")
	v2 := hertzer.NewRouterGroup("/v2")

	v2.Handle(&Ping{})
	v1.Handle(&Ping{})

	s.AddGroup(v1)
	v1.AddGroup(v2)

	if err := s.Run(); err != nil {
		panic(err)
	}
}

type Ping struct {
	httpx.MethodPost `route:"/ping/:name"`

	Name      string  `path:"name"`
	Age       int     `query:"age"`
	AuthToken string  `header:"AuthToken"`
	Address   Address `json:"address"`
	Score     int     `json:"score" form:"score"`
}

type Address struct {
	Home   string `json:"home" form:"home"`
	School string `json:"school"`
}

// func (Ping) Route() string {
// 	return "/ping/:name"
// }

// func (Ping) Method() string {
// 	return http.MethodGet
// }

func (p *Ping) Handle(ctx context.Context, arc *app.RequestContext) (any, error) {
	fmt.Println("handle ping")

	// (1) return response and nil error
	// return p, nil

	// (2) return response and error
	// err := fmt.Errorf("Origin Error")
	// return p, err

	// (3) return status response and status error
	// serr := errors.New(err, *p)
	// serr = serr.SetMessage("Error Message")
	// ret := resp.NewStatusResponse(consts.StatusBadGateway, *p)
	// return ret, serr

	// (4) return status response and nil error
	ret := resp.NewStatusResponse(consts.StatusAccepted, *p)
	return ret, nil
}

func (Ping) PreHandlers() []app.HandlerFunc {
	return []app.HandlerFunc{
		prefunc(),
	}
}

func (Ping) PostHandlers() []app.HandlerFunc {
	return []app.HandlerFunc{
		postfunc(),
	}
}

func prefunc() app.HandlerFunc {
	return func(ctx context.Context, arc *app.RequestContext) {
		fmt.Println("pre handler")
	}
}

func postfunc() app.HandlerFunc {
	return func(ctx context.Context, arc *app.RequestContext) {
		fmt.Println("post handler")
	}
}
