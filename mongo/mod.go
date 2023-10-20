package mongo

import (
	"github.com/qianwj/typed/mongo/client"
	"github.com/qianwj/typed/mongo/model"
)

func FromUri(uri string) *client.ClientBuilder {
	return client.NewClient().ApplyUri(uri)
}

func FromAddr(host string, port int) *client.ClientBuilder {
	addr := &model.Addr{
		Host: host,
		Port: port,
	}
	return FromAddrs(addr)
}

func FromAddrs(addrs ...*model.Addr) *client.ClientBuilder {
	hosts := make([]string, len(addrs))
	for i, addr := range addrs {
		hosts[i] = addr.String()
	}
	return client.NewClient().Hosts(hosts)
}
