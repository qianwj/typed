package microservice

import (
	"context"
	tfx "github.com/qianwj/typed/fx"
	data "github.com/qianwj/typed/fx/data"
	"log"
	"microservice/conf"
)

var components = []tfx.Component{}

func Bootstrap(options ...Option) {
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
		log.Fatal("init datasource error", err)
		components = append(components, dataSources...)
	}
}
