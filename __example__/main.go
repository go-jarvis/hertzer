package main

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-jarvis/herts"
	"github.com/go-jarvis/herts/pkg/httpx"
)

func main() {

	s := &herts.Server{
		Listen: ":8081",
	}

	s.Handle(&Ping{})

	v1 := herts.NewRouterGroup("/v1")
	v2 := herts.NewRouterGroup("/v2")

	v2.Handle(&Ping{})
	v1.Handle(&Ping{})

	s.AddGroup(v1)
	v1.AddGroup(v2)

	err := s.Run()
	if err != nil {
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
	return *p, nil
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
