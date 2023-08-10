package aggregates

import (
	"github.com/qianwj/typed/mongo/model"
	"github.com/qianwj/typed/mongo/model/filters"
	"github.com/qianwj/typed/mongo/operator"
	"github.com/qianwj/typed/mongo/util"
	"go.mongodb.org/mongo-driver/bson"
)

func (p *Pipeline) Count(field string) *Pipeline {
	if field == "" {
		field = "count"
	}
	p.put(operator.Count, field)
	return p
}

func (p *Pipeline) Facet(facets ...*model.Facet) *Pipeline {
	p.put(operator.Facet, bson.D(util.Map(facets, func(f *model.Facet) bson.E {
		return f.Marshal()
	})))
	return p
}

func (p *Pipeline) Match(query ...*filters.Filter) *Pipeline {
	p.put(operator.Match, util.Map(query, func(f *filters.Filter) bson.D {
		return f.Marshal()
	}))
	return p
}
