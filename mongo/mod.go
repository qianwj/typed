package mongo

import (
	"github.com/qianwj/typed/mongo/builder"
	"github.com/qianwj/typed/mongo/model"
)

func FromUri(uri string) *builder.ClientBuilder {
	return builder.NewClient().ApplyUri(uri)
}

func FromAddr(host string, port int) *builder.ClientBuilder {
	addr := &model.Addr{
		Host: host,
		Port: port,
	}
	return FromAddrs(addr)
}

func FromAddrs(addrs ...*model.Addr) *builder.ClientBuilder {
	hosts := make([]string, len(addrs))
	for i, addr := range addrs {
		hosts[i] = addr.String()
	}
	return builder.NewClient().Hosts(hosts)
}
