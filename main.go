package main

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
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
}

func (Ping) Route() string {
	return "/ping"
}

func (Ping) Method() string {
	return http.MethodGet
}

func (p *Ping) Handle(ctx context.Context, arc *app.RequestContext) (any, error) {
	return "pong", nil
}
