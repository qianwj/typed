package grpc

import (
	"context"
	tfx "github.com/qianwj/typed/fx"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type grpcClient struct {
	internal *grpc.ClientConn
}

func Apply(url string) (tfx.DataSource, error) {
	return nil, nil
}

func (g *grpcClient) Name(name string) tfx.DataSource {
	return g
}

func (g *grpcClient) Provide() fx.Option {
	return nil
}

func (g *grpcClient) Connect(ctx context.Context) error {
	return nil
}

func (g *grpcClient) Close(ctx context.Context) error {
	return nil
}
