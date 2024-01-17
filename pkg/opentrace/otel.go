package opentrace

import (
	"encoding/base64"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var (
	demoAuth = "demoUser:demoPass"
)

type Config struct {
	ServiceName string `env:""`
	Endpoint    string `env:""`
	BasicAuth   string `env:""`

	headers map[string]string
}

func (c *Config) SetName(name string) {
	c.ServiceName = name
}

func (c *Config) SetDefaults() {
	if c.ServiceName == "" {
		c.ServiceName = "echo"
	}

	if c.Endpoint == "" {
		c.Endpoint = "https://your-endpoint.com:4317"
	}

	if c.BasicAuth == "" {
		c.BasicAuth = demoAuth
	}

	if c.BasicAuth != "" && c.BasicAuth != demoAuth {
		if c.headers == nil {
			c.headers = map[string]string{}
		}
		c.headers["Authorization"] = BasicAuth(c.BasicAuth)
	}
}

func (c *Config) Initialize(opts ...Option) {
	c.SetDefaults()
	c.SetProvider()
}

func (c *Config) WithOptions(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}
}

func BasicAuth(auth string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

func (c *Config) SetProvider() trace.TracerProvider {
	tp, err := c.provider()
	if err != nil {
		return nil
	}

	otel.SetTracerProvider(tp)

	return tp
}

func GetProvider() trace.TracerProvider {
	return otel.GetTracerProvider()
}

type Option func(*Config)

func WithHeaders(headers map[string]string) Option {
	return func(c *Config) {
		if len(headers) == 0 {
			return
		}

		if len(c.headers) == 0 {
			c.headers = map[string]string{}
		}
		for k, v := range headers {
			c.headers[k] = v
		}
	}
}
