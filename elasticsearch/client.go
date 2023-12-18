package elasticsearch

import (
	"elasticsearch/api"
	"github.com/elastic/go-elasticsearch/v8"
)

type Client struct {
	internal *elasticsearch.TypedClient
}

func NewClient(builder ConfigBuilder) (*Client, error) {
	client, err := elasticsearch.NewTypedClient(builder.build())
	if err != nil {
		return nil, err
	}
	client.Search().Index("")
	return &Client{internal: client}, nil
}

func (c *Client) Indices() *api.IndicesAPI {
	return api.Indices(c.internal.Transport)
}

func (c *Client) Index(name string) *api.IndexAPI {
	return api.Index(name, c.internal.Transport)
}
