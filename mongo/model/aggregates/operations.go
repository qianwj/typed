package aggregates

import (
	"github.com/qianwj/typed/mongo/model/aggregates/group"
	"github.com/qianwj/typed/mongo/model/aggregates/lookup"
	"github.com/qianwj/typed/mongo/model/filters"
	"github.com/qianwj/typed/mongo/model/sorts"
	"github.com/qianwj/typed/mongo/operator"
	rawbson "go.mongodb.org/mongo-driver/bson"
)

func (p *Pipeline) Count(field string) *Pipeline {
	if field == "" {
		field = "count"
	}
	p.append(operator.Count, field)
	return p
}

func (p *Pipeline) Group(id group.ID, fields ...group.Accumulator) *Pipeline {
	body := rawbson.M{"_id": id}
	for _, acc := range fields {
		body[acc.Key] = acc.Value
	}
	p.append(operator.Group, body)
	return p
}

func (p *Pipeline) Lookup(cond *lookup.JoinCondition) *Pipeline {
	p.append(operator.Lookup, cond)
	return p
}

func (p *Pipeline) Match(filter *filters.Filter) *Pipeline {
	p.append(operator.Match, filter)
	return p
}

func (p *Pipeline) Skip(skip int64) *Pipeline {
	p.append(operator.Skip, skip)
	return p
}

func (p *Pipeline) Sort(opts *sorts.Options) *Pipeline {
	p.append(operator.Sort, opts)
	return p
}

func (p *Pipeline) SortByCount(expression any) *Pipeline {
	p.append(operator.SortByCount, expression)
	return p
}

func (p *Pipeline) Unset(fields ...string) *Pipeline {
	p.append(operator.Unset, fields)
	return p
}

//func (p *Pipeline) Facet(facets ...*model.Facet) *Pipeline {
//	p.put(operator.Facet, bson.D(util.Map(facets, func(f *model.Facet) primitive.E {
//		return f.Marshal()
//	})))
//	return p
//}
