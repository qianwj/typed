package grpc

import (
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/qianwj/typed/fx/options"
	"google.golang.org/grpc"
)

const (
	maxConcurrentStreamsKey = "maxConcurrentStreams"
	validateKey             = "validate"
	recoveryKey             = "recovery"
)

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

func (o *grpcOptions) GetValue(key string) (any, bool) {
	switch key {
	case maxConcurrentStreamsKey:
		return o.maxConcurrentStreams, o.maxConcurrentStreams > 0
	case validateKey:
		return o.validate, o.validate
	case recoveryKey:
		return o.recovery, o.recovery
	}
	return o, false
}

func MaxConcurrentStreams(maxConcurrentStreams uint32) *grpcOptions {
	return &grpcOptions{maxConcurrentStreams: maxConcurrentStreams}
}

func EnableValidate() *grpcOptions {
	return &grpcOptions{validate: true}
}

func mergeOptions(opts ...options.Options) *grpcOptions {
	result := new(grpcOptions)
	for _, opt := range opts {
		addr, ok := opt.GetValue(options.AddressKey)
		if ok {
			result.addr = addr.(string)
			continue
		}
		maxConcurrentStreams, ok := opt.GetValue(maxConcurrentStreamsKey)
		if ok {
			result.maxConcurrentStreams = maxConcurrentStreams.(uint32)
		}
		if _, ok = opt.GetValue(validateKey); ok {
			result.validate = true
		}
		if _, ok = opt.GetValue(recoveryKey); ok {
			result.recovery = true
		}
		metricsPort, _ := opt.GetValue(options.MetricsPortKey)
		if ok {
			result.metrics = &struct{ port int }{port: metricsPort.(int)}
		}
	}
	return result
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
