package main

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-jarvis/herts/pkg/httpx"
	"github.com/go-jarvis/herts/server"
)

func main() {

	s := &server.Server{
		Listen: ":8081",
	}

	s.Handle(&Ping{})

	v1 := server.NewRouterGroup("/v1")
	v2 := server.NewRouterGroup("/v2")

	v2.Handle(&Ping{})
	v1.Handle(&Ping{})

	s.AppendGroup(v1)
	v1.AppendGroup(v2)

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
	return *p, nil
}
