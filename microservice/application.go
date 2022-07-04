package microservice

import (
	tfx "github.com/qianwj/typed/fx"
	"log"
	"microservice/conf"
)

var components = []tfx.Component{}

func Bootstrap(options ...Option) {
	appConf := defaultConf()
	for _, option := range options {
		option(appConf)
	}
	err := conf.Load("")
	if err != nil {
		log.Fatal("reading conf error", err)
	}
	if appConf.enableData {

	}
}
