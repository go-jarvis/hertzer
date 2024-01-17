package opentrace

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

// Initializes an OTLP exporter, and configures the corresponding trace and
// metric providers.
func (c *Config) provider() (trace.TracerProvider, error) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceNameKey.String(c.ServiceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// If the OpenTelemetry Collector is running on a local cluster (minikube or
	// microk8s), it should be accessible through the NodePort service at the
	// `localhost:30080` endpoint. Otherwise, replace `localhost` with the
	// endpoint of your cluster. If you run the app inside k8s, then you can
	// probably connect directly to the service through dns.
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	// var err error
	exporter, err := c.traceExporter(ctx, c.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := sdktrace.NewBatchSpanProcessor(exporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	return tracerProvider, nil
}

func (c *Config) traceExporter(ctx context.Context, endpoint string) (*otlptrace.Exporter, error) {
	ur, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	var exporter *otlptrace.Exporter

	switch strings.ToLower(ur.Scheme) {
	case "http", "https":
		exporter, err = c.httpExporter(ctx, ur)
		if err != nil {
			return nil, err
		}
	case "grpc":
		fallthrough
	default:
		exporter, err = c.grpcExpoter(ctx, ur)
		if err != nil {
			return nil, err
		}
	}

	return exporter, nil
}

func (c *Config) httpExporter(ctx context.Context, ur *url.URL) (*otlptrace.Exporter, error) {

	opts := []otlptracehttp.Option{
		otlptracehttp.WithTimeout(5 * time.Second),
		otlptracehttp.WithEndpoint(ur.Host),
		// otlptracehttp.WithInsecure(),
		otlptracehttp.WithHeaders(
			c.headers,
		),
	}

	if ur.Scheme != "https" {
		opts = append(opts, otlptracehttp.WithInsecure())
	}

	trace, err := otlptracehttp.New(ctx, opts...)

	return trace, err
}

// 创建 OTEL 的 GRPC 连接器
func (c *Config) grpcExpoter(ctx context.Context, ur *url.URL) (*otlptrace.Exporter, error) {
	// addr := strings.TrimLeft(endpoint, "grpc://")

	conn, err := grpc.DialContext(ctx, ur.Host,
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		// grpc.WithTimeout(5*time.Second),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	// Set up a trace exporter
	traceExporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithGRPCConn(conn),
		otlptracegrpc.WithHeaders(
			c.headers,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}
	return traceExporter, nil
}
