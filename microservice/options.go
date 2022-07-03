package microservice

import (
	"github.com/qianwj/typed/collection/set"
)

type serverType int
type loggerType int

const (
	serverGrpc serverType = iota + 1
	serverHttp
	loggerZap  = iota
	dataMongo  = "mongo"
	dataRedis  = "redis"
	cacheRedis = "redis"
)

type appConf struct {
	serverType  serverType
	loggerType  loggerType
	enableData  bool
	dataTypes   *set.Set[string]
	enableCache bool
	cacheType   string
}

type Option func(*appConf)

func defaultConf() *appConf {
	return &appConf{
		serverType:  serverGrpc,
		enableData:  false,
		dataTypes:   set.NewSet[string](),
		enableCache: false,
		cacheType:   "",
	}
}

func Grpc() Option {
	return func(c *appConf) {
		c.serverType = serverGrpc
	}
}

func Http() Option {
	return func(c *appConf) {
		c.serverType = serverHttp
	}
}

func Zap() Option {
	return func(c *appConf) {
		c.loggerType = loggerZap
	}
}

func Mongo() Option {
	return func(c *appConf) {
		c.enableData = true
		c.dataTypes.Add(dataMongo)
	}
}

func Redis() Option {
	return func(c *appConf) {
		c.enableData = true
		c.dataTypes.Add(dataRedis)
	}
}
