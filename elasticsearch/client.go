package elasticsearch

import (
	"elasticsearch/api"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/ping"
)

type Client struct {
	internal *elasticsearch.TypedClient
}

func NewClient(builder ConfigBuilder) (*Client, error) {
	client, err := elasticsearch.NewTypedClient(builder.build())
	if err != nil {
		return nil, err
	}
	return &Client{internal: client}, nil
}

func (c *Client) Ping() *ping.Ping {
	return c.internal.Ping()
}

func (c *Client) Indices() *api.IndicesAPI {
	return api.Indices(c.internal.Transport)
}

func (c *Client) Index(name string) *api.IndexAPI {
	return api.Index(name, c.internal.Transport)
}
