package server

import (
	"context"
	"github.com/qianwj/typed/config"
	tfx "github.com/qianwj/typed/fx"
	"github.com/qianwj/typed/fx/server/grpc"
)

type Config struct {
	Grpc *grpc.Config
}

func Bootstrap(ctx context.Context) ([]tfx.Component, error) {
	components := make([]tfx.Component, 0)
	conf, err := config.Unmarshal[Config]("server")
	if err != nil {
		return nil, err
	}
	if conf.Grpc != nil {
		server := grpc.NewServer(grpc.FromConfig(conf.Grpc))
		components = append(components, server)
	}
	return components, nil
}
