package redis

import (
	"context"
	"fmt"
	client "github.com/go-redis/redis"
	tfx "github.com/qianwj/typed/fx"
	"go.uber.org/fx"
)

type redisClient struct {
	name     string
	internal *client.Client
}

func NewData(opt *client.Options) (tfx.DataAccess, error) {
	c := client.NewClient(opt)
	return &redisClient{internal: c}, nil
}

func Apply(uri string) (tfx.DataAccess, error) {
	opt, err := client.ParseURL(uri)
	if err != nil {
		return nil, err
	}
	return &redisClient{
		internal: client.NewClient(opt),
	}, nil
}

func (r *redisClient) Name(name string) tfx.DataAccess {
	r.name = name
	return r
}

func (r *redisClient) Provide() fx.Option {
	data := fx.Provide(fx.Annotate(func() tfx.DataAccess {
		return r
	}), fx.ResultTags(`group:"data_sources"`))
	tag := `name:"redis"`
	if len(r.name) > 0 {
		tag = fmt.Sprintf(`name:"%s_redis"`, r.name)
	}
	return fx.Options(data, fx.Provide(fx.Annotate(r.client, fx.ResultTags(tag))))
}

func (r *redisClient) client() *client.Client {
	return r.internal
}

func (r *redisClient) Connect(ctx context.Context) error {
	r.internal.WithContext(ctx)
	res := r.internal.Ping()
	return res.Err()
}

func (r *redisClient) Close(ctx context.Context) error {
	return r.internal.Close()
}
