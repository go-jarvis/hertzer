package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/go-jarvis/hertzer"
	"github.com/go-jarvis/hertzer/pkg/common/middlewares/logr"
	"github.com/go-jarvis/slogr"
)

func main() {

	// err := errors.ErrBadPoolConn

	s := &hertzer.Server{
		Listen: ":8081",
	}

	log := slogr.Default()
	logc := logr.AcceccLoggerConfig{
		SkipPaths: []string{
			"/api/ping/wangwu",
		},
	}

	s.Use(
		logr.WithSLogger(log),
		logr.AccessLoggerWithConfig(logc),
	)

	// 注册路由
	register(s)

	if err := s.Run(); err != nil {
		panic(err)
	}
}

// register 注册路由
func register(s *hertzer.Server) {

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

}
