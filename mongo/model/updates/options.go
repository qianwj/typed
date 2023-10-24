package updates

import (
	"github.com/qianwj/typed/mongo/model/filters"
	"github.com/qianwj/typed/mongo/util"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ArrayFilters struct {
	items    []*filters.Filter
	registry *bsoncodec.Registry
}

func NewArrayFilters(filters ...*filters.Filter) *ArrayFilters {
	return &ArrayFilters{
		items: filters,
	}
}

func (af *ArrayFilters) Registry(r *bsoncodec.Registry) *ArrayFilters {
	af.registry = r
	return af
}

func (af *ArrayFilters) Raw() options.ArrayFilters {
	res := options.ArrayFilters{
		Filters: util.ToAny(af.items),
	}
	if af.registry != nil {
		// Registry is the registry to use for converting filters. Defaults to bson.DefaultRegistry.
		//
		// Deprecated: Marshaling ArrayFilters to BSON will not be supported in Go Driver 2.0.
		res.Registry = af.registry
	}
	return res
}
