package server

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/go-jarvis/herts/route"
)

type Config struct {
	Listen string `env:""`

	h *server.Hertz
	r *route.RouterGroup

	opts []config.Option
}

func (s *Config) SetDefaults() {
	if s.Listen == "" {
		s.Listen = ":8080"
	}

	hp := server.WithHostPorts(s.Listen)
	s.WithOptions(hp)
}

func (s *Config) initialize() {
	s.h = server.Default(s.opts...)
}

func (s *Config) Run() error {
	s.SetDefaults()
	s.initialize()

	return s.h.Run()
}

func (s *Config) WithOptions(opts ...config.Option) {
	if len(s.opts) == 0 {
		s.opts = make([]config.Option, 0)
	}

	s.opts = append(s.opts, opts...)
}
