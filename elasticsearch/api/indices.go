package api

import (
	"elasticsearch/requests"
	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/mget"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/close"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/delete"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/exists"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/get"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/open"
)

type IndicesAPI struct {
	tp elastictransport.Interface
}

func Indices(tp elastictransport.Interface) *IndicesAPI {
	return &IndicesAPI{tp: tp}
}

func (i *IndicesAPI) Create(index string, req *requests.CreateIndexRequestBuilder) *create.Create {
	return create.NewCreateFunc(i.tp)(index).Request(req.Build())
}

func (i *IndicesAPI) Delete(index string) *delete.Delete {
	return delete.NewDeleteFunc(i.tp)(index)
}

func (i *IndicesAPI) Exists(index string) *exists.Exists {
	return exists.NewExistsFunc(i.tp)(index)
}

func (i *IndicesAPI) Close(index string) *close.Close {
	return close.NewCloseFunc(i.tp)(index)
}

func (i *IndicesAPI) Open(index string) *open.Open {
	return open.NewOpenFunc(i.tp)(index)
}

func (i *IndicesAPI) Get(index string) *get.Get {
	return get.NewGetFunc(i.tp)(index)
}

func (i *IndicesAPI) MGet() *mget.Mget {
	return mget.NewMgetFunc(i.tp)()
}
