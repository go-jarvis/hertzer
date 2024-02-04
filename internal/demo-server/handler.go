package main

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-jarvis/hertzer/pkg/httpx"
	"github.com/go-jarvis/slogr"
)

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

func (p *Ping) Handle(ctx context.Context, arc *app.RequestContext) {
	// fmt.Println("handle ping")
	log := slogr.FromContext(ctx)
	log.Info("handle ping")

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
	// ret := resp.NewStatusResponse(consts.StatusAccepted, *p)
	// return ret, nil
	// arc.JSON(200, ret)
	arc.JSON(200, p)
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
