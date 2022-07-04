package grpc

import (
	"context"
	"fmt"
	"github.com/Masterminds/log-go"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	tfx "github.com/qianwj/typed/fx"
	"go.uber.org/atomic"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
)

type Server struct {
	opts *grpcOptions
	srv  *grpc.Server
	ctx  *applicationContext
}

type applicationContext struct {
	services []fx.Option
	state    *atomic.Int32
}

func (ctx *applicationContext) provide() fx.Option {
	return fx.Options(
		fx.Options(ctx.services...),
	)
}

func NewServer(options ...Option) tfx.Server {
	opts := new(grpcOptions)
	for _, option := range options {
		option(opts)
	}
	return &Server{opts: opts, ctx: &applicationContext{}}
}

func (app *Server) Provide() fx.Option {
	return fx.Options(app.ctx.provide(), fx.Supply(app), fx.Invoke(runGrpcServer))
}

func (app *Server) onStart(ctx context.Context) error {
	lis, err := net.Listen("tcp", app.opts.addr)
	if err != nil {
		return err
	}
	go func(ctx context.Context) {
		if err := app.srv.Serve(lis); err != nil {
			log.Error("starting grpc server error:", err)
			os.Exit(-1)
		}
	}(ctx)
	return nil
}

func (app *Server) onStop(ctx context.Context) error {
	app.srv.GracefulStop()
	return nil
}

func runGrpcServer(app *Server, serviceModule serviceModule, lifecycle fx.Lifecycle) *grpc.Server {
	srv := grpc.NewServer(grpcServerOptions(app.opts)...)
	serviceModule.register(srv)
	app.srv = srv
	if app.opts.metrics != nil {
		exportMetrics(app.opts.metrics.port, srv)
	}
	lifecycle.Append(fx.Hook{
		OnStart: app.onStart,
		OnStop:  app.onStop,
	})
	return srv
}

func exportMetrics(port int, srv *grpc.Server) {
	registry := prometheus.NewRegistry()
	registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	registry.MustRegister(collectors.NewGoCollector())
	registry.MustRegister(grpc_prometheus.DefaultServerMetrics)
	grpc_prometheus.Register(srv)
	server := &http.Server{Handler: promhttp.HandlerFor(registry, promhttp.HandlerOpts{}), Addr: fmt.Sprintf(":%d", port)}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			zap.L().Error("exporter start error: ", zap.Error(err))
		}
	}()
}
