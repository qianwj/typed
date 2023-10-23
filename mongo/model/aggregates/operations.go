package aggregates

import (
	"github.com/qianwj/typed/mongo/model/filters"
	"github.com/qianwj/typed/mongo/model/sorts"
	"github.com/qianwj/typed/mongo/operator"
)

func (p *Pipeline) Count(field string) *Pipeline {
	if field == "" {
		field = "count"
	}
	p.put(operator.Count, field)
	return p
}

func (p *Pipeline) Match(filter *filters.Filter) *Pipeline {
	p.put(operator.Match, filter)
	return p
}

func (p *Pipeline) Skip(skip int64) *Pipeline {
	p.put(operator.Skip, skip)
	return p
}

func (p *Pipeline) Sort(opts *sorts.Options) *Pipeline {
	p.put(operator.Sort, opts)
	return p
}

func (p *Pipeline) SortByCount(expression any) *Pipeline {
	p.put(operator.SortByCount, expression)
	return p
}

func (p *Pipeline) Unset(fields ...string) *Pipeline {
	p.put(operator.Unset, fields)
	return p
}

//func (p *Pipeline) Facet(facets ...*model.Facet) *Pipeline {
//	p.put(operator.Facet, bson.D(util.Map(facets, func(f *model.Facet) primitive.E {
//		return f.Marshal()
//	})))
//	return p
//}
//
//func (p *Pipeline) Group(id GroupId, fields ...*GroupField) *Pipeline {
//	op := append(primitive.D{}, primitive.E{Key: "_id", Value: id})
//	op = append(op, util.Map(fields, func(f *GroupField) primitive.E {
//		return f.Marshal()
//	})...)
//	p.put(operator.Group, op)
//	return p
//}
//
//func (p *Pipeline) Lookup(l *Lookup) *Pipeline {
//	p.put(operator.Lookup, l.Marshal())
//	return p
//}
