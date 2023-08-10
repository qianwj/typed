package pipe

import (
	"github.com/qianwj/typed/mongo/model"
	"github.com/qianwj/typed/mongo/model/aggregate"
	"github.com/qianwj/typed/mongo/model/filters"
	"github.com/qianwj/typed/mongo/model/sorts"
	"github.com/qianwj/typed/mongo/operator"
	"go.mongodb.org/mongo-driver/bson"
)

func (p *Pipeline) Count(target string) *Pipeline {
	p.put(operator.Count, target)
	return p
}

func (p *Pipeline) Documents(docs []any) *Pipeline {
	p.put(operator.Documents, docs)
	return p
}

func (p *Pipeline) GraphLookup(join *aggregate.GraphLookup) *Pipeline {
	p.put(operator.Lookup, join.Marshal())
	return p
}

func (p *Pipeline) Group(id groupId, fields ...model.Pair[bson.M]) *Pipeline {
	val := bson.D{{Key: "_id", Value: id}}
	for _, field := range fields {
		val = append(val, bson.E{Key: field.Key, Value: field.Value})
	}
	p.put(operator.Group, val)
	return p
}

func (p *Pipeline) Limit(value int64) *Pipeline {
	p.put(operator.Limit, value)
	return p
}

func (p *Pipeline) Lookup(join *aggregate.Lookup) *Pipeline {
	p.put(operator.Lookup, join.Marshal())
	return p
}

func (p *Pipeline) Match(match *filters.Filter) *Pipeline {
	p.put(operator.Match, match.Marshal())
	return p
}

func (p *Pipeline) Project(fields ...model.Pair[bool]) *Pipeline {
	val := bson.D{}
	for _, field := range fields {
		val = append(val, bson.E{Key: field.Key, Value: field.Value})
	}
	p.put(operator.Project, val)
	return p
}

func (p *Pipeline) Sort(sort *sorts.SortOptions) *Pipeline {
	p.put(operator.Sort, sort.Marshal())
	return p
}

func (p *Pipeline) Unset(fields ...string) *Pipeline {
	p.put(operator.Unset, fields)
	return p
}
