package redis

import (
	"context"
	"fmt"
	client "github.com/go-redis/redis"
	tfx "github.com/qianwj/typed/fx"
	"go.uber.org/fx"
)

type redisClient struct {
	internal *client.Client
}

func NewData(opt *client.Options) (tfx.DataAccess, error) {
	c := client.NewClient(opt)
	return &redisClient{c}, nil
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

func (r *redisClient) Provide(name ...string) fx.Option {
	if len(name) == 0 {
		return fx.Provide(fx.Annotate(func() *client.Client {
			return r.internal
		}, fx.ResultTags(`name:"redis"`)))
	} else {
		return fx.Provide(fx.Annotate(func() *client.Client {
			return r.internal
		}, fx.ResultTags(fmt.Sprintf(`name:"%s_redis"`, name[0]))))
	}
}

func (r *redisClient) Connect(ctx context.Context) error {
	r.internal.WithContext(ctx)
	res := r.internal.Ping()
	return res.Err()
}

func (r *redisClient) Close(ctx context.Context) error {
	return r.internal.Close()
}
