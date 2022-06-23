package main

import (
	"context"
	"fx-grpc/api"
	fx_grpc "github.com/qianwj/typed/fx/grpc"
	"google.golang.org/grpc"
)

func main() {
	app := fx_grpc.NewApp()
	app.RegisterService(NewGreeterService)
	app.Run()
}

type GreeterService struct {
	api.UnimplementedGreeterServer
}

func (svc *GreeterService) Register(srv *grpc.Server) {
	api.RegisterGreeterServer(srv, svc)
}

func (svc *GreeterService) Hello(ctx context.Context, req *api.HelloRequest) (*api.HelloReply, error) {
	return &api.HelloReply{Message: "hello, " + req.GetName()}, nil
}

func NewGreeterService() fx_grpc.Service {
	return &GreeterService{}
}
