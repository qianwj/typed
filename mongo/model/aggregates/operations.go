package aggregates

import (
	"github.com/qianwj/typed/mongo/bson"
	"github.com/qianwj/typed/mongo/model"
	"github.com/qianwj/typed/mongo/model/filters"
	"github.com/qianwj/typed/mongo/operator"
	"github.com/qianwj/typed/mongo/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (p *Pipeline) Count(field string) *Pipeline {
	if field == "" {
		field = "count"
	}
	p.put(operator.Count, field)
	return p
}

func (p *Pipeline) Facet(facets ...*model.Facet) *Pipeline {
	p.put(operator.Facet, bson.D(util.Map(facets, func(f *model.Facet) primitive.E {
		return f.Marshal()
	})))
	return p
}

func (p *Pipeline) Group(id GroupId, pipeline bson.A) *Pipeline {
	p.put(operator.Group, bson.D{
		{Key: "_id", Value: id},
		//pipeline...,
	})
	return p
}

func (p *Pipeline) Match(query ...*filters.Filter) *Pipeline {
	p.put(operator.Match, util.Map(query, func(f *filters.Filter) primitive.D {
		return f.Marshal()
	}))
	return p
}
