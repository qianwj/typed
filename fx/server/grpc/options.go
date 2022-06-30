package grpc

import (
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
)

type Config struct {
	Port                 int
	EnableHealthCheck    bool
	EnableValidate       bool
	EnableRecovery       bool
	MaxConcurrentStreams uint32
	Metrics              *struct {
		Port int
	}
}

type grpcOptions struct {
	addr                 string
	healthCheck          bool
	validate             bool
	recovery             bool
	maxConcurrentStreams uint32
	metrics              *struct {
		port int
	}
	tracing bool
}

type Option func(o *grpcOptions)

func Address(addr string) Option {
	return func(o *grpcOptions) {
		o.addr = addr
	}
}

func MaxConcurrentStreams(maxConcurrentStreams uint32) Option {
	return func(o *grpcOptions) {
		o.maxConcurrentStreams = maxConcurrentStreams
	}
}

func EnableValidate() Option {
	return func(o *grpcOptions) {
		o.validate = true
	}
}

func EnableHealthCheck() Option {
	return func(o *grpcOptions) {
		o.healthCheck = true
	}
}

func EnableRecovery() Option {
	return func(o *grpcOptions) {
		o.recovery = true
	}
}

func Metrics(port int) Option {
	return func(o *grpcOptions) {
		o.metrics = &struct{ port int }{port: port}
	}
}

func FromConfig(c *Config) Option {
	return func(o *grpcOptions) {
		o.addr = fmt.Sprintf(":%d", c.Port)
		o.healthCheck = c.EnableHealthCheck
		o.validate = c.EnableValidate
		o.recovery = c.EnableRecovery
		o.maxConcurrentStreams = c.MaxConcurrentStreams
		if c.Metrics != nil {
			o.metrics = &struct{ port int }{port: c.Port}
		}
	}
}

func grpcServerOptions(opts *grpcOptions) []grpc.ServerOption {
	var unaryServerInterceptors []grpc.UnaryServerInterceptor
	var streamServerInterceptors []grpc.StreamServerInterceptor
	if opts.validate {
		unaryServerInterceptors = append(unaryServerInterceptors, grpc_validator.UnaryServerInterceptor())
		streamServerInterceptors = append(streamServerInterceptors, grpc_validator.StreamServerInterceptor())
	}
	if opts.recovery {
		unaryServerInterceptors = append(unaryServerInterceptors, grpc_recovery.UnaryServerInterceptor())
		streamServerInterceptors = append(streamServerInterceptors, grpc_recovery.StreamServerInterceptor())
	}
	if opts.metrics != nil {
		unaryServerInterceptors = append(unaryServerInterceptors, grpc_prometheus.UnaryServerInterceptor)
		streamServerInterceptors = append(streamServerInterceptors, grpc_prometheus.StreamServerInterceptor)
	}
	serverOptions := []grpc.ServerOption{
		grpc_middleware.WithUnaryServerChain(unaryServerInterceptors...),
		grpc_middleware.WithStreamServerChain(streamServerInterceptors...),
	}
	if opts.maxConcurrentStreams > 0 {
		serverOptions = append(serverOptions, grpc.MaxConcurrentStreams(opts.maxConcurrentStreams))
	}
	return serverOptions
}
