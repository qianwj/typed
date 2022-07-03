package microservice

import (
	tfx "github.com/qianwj/typed/fx"
	"github.com/qianwj/typed/fx/data/mongo"
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
		dataTypes := appConf.dataTypes
		for _, dataType := range dataTypes.Slice() {
			switch dataType {
			case dataMongo:
				mongoConf, err := conf.Unmarshal[map[string]mongo.Conf]("data.mongo")
				if err != nil {
					log.Fatal("reading mongo conf error:", err)
				}
				for name, config := range mongoConf {
					if name == "default" {
						component, err := mongo.Apply(config.Uri)
						if err != nil {
							log.Fatal("init mongo error:", err)
						}
						components = append(components, component)
					}
				}
			}
		}
	}
}
