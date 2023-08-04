package options

import (
	"github.com/qianwj/typed/mongo/model/filter"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ArrayFilters struct {
	Items    []*filter.Filter
	Registry *bsoncodec.Registry
}

func (af *ArrayFilters) Raw() options.ArrayFilters {
	// complete this function
	return options.ArrayFilters{}
}
