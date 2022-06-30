package microservice

import (
	"go.uber.org/fx"
	"log"
	"microservice/conf"
)

var beans = []fx.Option{}

func Bootstrap(options ...Option) {
	appConf := defaultConf()
	for _, option := range options {
		option(appConf)
	}
	err := conf.Load("")
	if err != nil {
		log.Fatal("reading configuration error", err)
	}

}
