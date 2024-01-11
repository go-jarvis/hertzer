package hertzer

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/go-jarvis/hertzer/pkg/common/middlewares/monitor"
	"github.com/go-jarvis/hertzer/pkg/operator"
)

type Server struct {
	Listen string `env:""`

	h *server.Hertz
	r *RouterGroup

	opts []config.Option
}

func (s *Server) SetDefaults() {
	if s.Listen == "" {
		s.Listen = ":8080"
	}

	hp := server.WithHostPorts(s.Listen)
	s.WithOptions(
		hp,
	)
}

func (s *Server) defaultRouterGroup() {
	if s.r == nil {
		s.r = NewRouterGroup("/")
	}
}

func (s *Server) initialize() {

	s.h = server.Default(s.opts...)
	s.Handle(&monitor.Liveness{})

	s.defaultRouterGroup()
	s.r.r = s.h.Group("/")

	s.r.initialize()
}

func (s *Server) Run() error {
	s.SetDefaults()
	s.initialize()

	return s.h.Run()
}

func (s *Server) WithOptions(opts ...config.Option) {
	if len(s.opts) == 0 {
		s.opts = make([]config.Option, 0)
	}

	s.opts = append(s.opts, opts...)
}

func (s *Server) Use(middleware ...app.HandlerFunc) {
	s.defaultRouterGroup()

	s.r.Use(middleware...)
}

func (s *Server) Handle(opers ...operator.Operator) {
	s.defaultRouterGroup()

	s.r.Handle(opers...)
}

func (s *Server) AddGroup(group ...*RouterGroup) {
	s.defaultRouterGroup()

	s.r.AddGroup(group...)
}
