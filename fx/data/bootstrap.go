package data

import (
	"context"
	"github.com/Masterminds/log-go"
	"github.com/qianwj/typed/collection/set"
	"github.com/qianwj/typed/config"
	tfx "github.com/qianwj/typed/fx"
	"github.com/qianwj/typed/fx/data/mongo"
	"github.com/qianwj/typed/fx/data/redis"
	"golang.org/x/sync/errgroup"
)

const (
	TypeMongo = "mongo"
	TypeRedis = "redis"
)

func Bootstrap(ctx context.Context, dataTypes set.Set[string]) ([]tfx.Component, error) {
	components := make([]tfx.Component, 0)
	var group errgroup.Group
	for _, dataType := range dataTypes.Slice() {
		switch dataType {
		case TypeMongo:
			mongoConf, err := config.Unmarshal[map[string]mongo.Conf]("data.mongo")
			if err != nil {
				log.Fatal("reading mongo conf error:", err)
			}
			for name, conf := range *mongoConf {
				component, err := mongo.Apply(conf.Uri)
				if err != nil {
					log.Fatal("init mongo error:", err)
				}
				if name != "default" {
					component.Name(name)
				}
				group.Go(func() error {
					return component.Connect(ctx)
				})
				components = append(components, component)
			}
		case TypeRedis:
			redisConf, err := config.Unmarshal[map[string]redis.Conf]("data.redis")
			if err != nil {
				log.Fatal("reading redis conf error:", err)
			}
			for name, conf := range *redisConf {
				component, err := redis.Apply(conf.Uri)
				if err != nil {
					log.Fatal("init redis error:", err)
				}
				if name != "default" {
					component.Name(name)
				}
				group.Go(func() error {
					return component.Connect(ctx)
				})
				components = append(components, component)
			}
		}
	}
	if err := group.Wait(); err != nil {
		return nil, err
	}
	return components, nil
}
