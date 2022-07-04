package microservice

import (
	"context"
	tfx "github.com/qianwj/typed/fx"
	"github.com/qianwj/typed/fx/data"
	"github.com/qianwj/typed/fx/server"
	"go.uber.org/fx"
	"log"
	"microservice/conf"
)

var (
	app        *tfx.Application
	components = make([]tfx.Component, 0)
)

func Bootstrap(options ...Option) *tfx.Application {
	ctx := context.TODO()
	appConf := defaultConf()
	for _, option := range options {
		option(appConf)
	}
	err := conf.Load("")
	if err != nil {
		log.Fatal("reading conf error", err)
	}
	if appConf.enableData {
		dataSources, err := data.Bootstrap(ctx, *appConf.dataTypes)
		if err != nil {
			log.Fatal("init datasource error", err)
		}
		components = append(components, dataSources...)
	}
	servers, err := server.Bootstrap(ctx)
	if err != nil {
		log.Fatal("init server error", err)
	}
	components = append(components, servers...)
	app = tfx.NewApp(components...)
	return app
}

func Provide(components ...fx.Option) {
	app.Provide(components...)
}

func Run() {
	app.Run()
}
