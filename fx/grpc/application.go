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
	"github.com/qianwj/typed/fx/options"
	"github.com/qianwj/typed/fx/util"
	"go.uber.org/atomic"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"net/http"
)

const (
	stateInit = iota
	stateRunning
)

type Application struct {
	services []fx.Option
	opts     *grpcOptions
	srv      *grpc.Server
	state    *atomic.Int32
}

func NewApp(options ...options.Options) tfx.Application {
	return &Application{opts: mergeOptions(options...), state: atomic.NewInt32(stateInit)}
}

func (app *Application) WithLogger(constructor any) tfx.Application {
	return app
}

func (app *Application) Run() {
	app.state.Inc()
	container := fx.New(fx.Module("service", app.services...), fx.Supply(app), fx.Invoke(runGrpcServer))
	go func() {
		signal := <-container.Done()
		log.Debugf("receive signal(%s), app shutdown", signal.String())
	}()
	container.Run()
}

func (app *Application) onStart(ctx context.Context) error {
	lis, err := net.Listen("tcp", app.opts.addr)
	if err != nil {
		return err
	}
	go func() {
		if err := app.srv.Serve(lis); err != nil {
			util.Panic(err)
		}
	}()
	return nil
}

func (app *Application) onStop(ctx context.Context) error {
	app.srv.GracefulStop()
	return nil
}

func runGrpcServer(app *Application, serviceModule serviceModule, lifecycle fx.Lifecycle) *grpc.Server {
	srv := grpc.NewServer(grpcServerOptions(app.opts)...)
	serviceModule.register(srv)
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
