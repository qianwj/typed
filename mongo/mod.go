package mongo

import (
	"github.com/qianwj/typed/mongo/client"
	"github.com/qianwj/typed/mongo/model"
)

func FromUri(uri string) *client.Builder {
	return client.NewBuilder().ApplyUri(uri)
}

func FromAddr(host string, port int) *client.Builder {
	addr := &model.Addr{
		Host: host,
		Port: port,
	}
	return FromAddrs(addr)
}

func FromAddrs(addrs ...*model.Addr) *client.Builder {
	hosts := make([]string, len(addrs))
	for i, addr := range addrs {
		hosts[i] = addr.String()
	}
	return client.NewBuilder().Hosts(hosts)
}
