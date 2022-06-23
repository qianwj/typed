package grpc

import (
	"github.com/Masterminds/log-go"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type Service interface {
	Register(srv *grpc.Server)
}

func (app *Application) RegisterService(constructors ...any) {
	for _, constructor := range constructors {
		opt := fx.Provide(fx.Annotated{Target: constructor, Group: "svc"})
		app.services = append(app.services, opt)
	}
}

type serviceModule struct {
	fx.In
	Services []Service
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
