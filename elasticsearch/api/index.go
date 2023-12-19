package api

import (
	"elasticsearch/requests"
	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/bulk"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/count"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/create"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/delete"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/explain"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/get"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/index"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/update"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

type IndexAPI struct {
	index string
	tp    elastictransport.Interface
}

func Index(index string, tp elastictransport.Interface) *IndexAPI {
	return &IndexAPI{index: index, tp: tp}
}

func (i *IndexAPI) Create(id string, doc any) *create.Create {
	return create.NewCreateFunc(i.tp)(i.index, id).Request(doc)
}

func (i *IndexAPI) UpdateById(id string, req *requests.UpdateRequestBuilder) *update.Update {
	return update.NewUpdateFunc(i.tp)(i.index, id).Request(req.Build())
}

func (i *IndexAPI) Delete(id string) *delete.Delete {
	return delete.NewDeleteFunc(i.tp)(i.index, id)
}

func (i *IndexAPI) Indexing(id string) *index.Index {
	return index.NewIndexFunc(i.tp)(i.index).Id(id)
}

func (i *IndexAPI) Bulk(docs []any) *bulk.Bulk {
	req := docs
	return bulk.NewBulkFunc(i.tp)().Index(i.index).Request(&req)
}

func (i *IndexAPI) Get(id string) *get.Get {
	return get.NewGetFunc(i.tp)(i.index, id)
}

func (i *IndexAPI) Search(req *requests.SearchRequestBuilder) *search.Search {
	return search.NewSearchFunc(i.tp)().Index(i.index).Request(req.Build())
}

func (i *IndexAPI) Count(query *types.Query) *count.Count {
	req := &count.Request{Query: query}
	return count.NewCountFunc(i.tp)().Index(i.index).Request(req)
}

func (i *IndexAPI) Explain(id string, query *types.Query) *explain.Explain {
	req := &explain.Request{Query: query}
	return explain.NewExplainFunc(i.tp)(i.index, id).Request(req)
}
