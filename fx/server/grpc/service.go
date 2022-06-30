package grpc

import (
	"github.com/Masterminds/log-go"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type Service interface {
	Register(srv *grpc.Server)
}

func (app *Server) RegisterService(constructors ...any) {
	ctx := app.ctx
	for _, constructor := range constructors {
		opt := fx.Provide(fx.Annotate(constructor, fx.ResultTags(`group:"grpc_service"`)))
		ctx.services = append(ctx.services, opt)
	}
	app.ctx = ctx
}

type serviceModule struct {
	fx.In
	Services []Service `group:"grpc_service"`
}

func (s *serviceModule) register(srv *grpc.Server) {
	if len(s.Services) == 0 {
		log.Warn("no grpc service! please use `app.RegisterService()` to register service")
	}
	log.Debugf("register(%d) services...", len(s.Services))
	for _, service := range s.Services {
		service.Register(srv)
	}
}
