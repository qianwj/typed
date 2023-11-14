package aggregates

import (
	"github.com/qianwj/typed/mongo/bson"
	"github.com/qianwj/typed/mongo/model/aggregates/lookup"
	"github.com/qianwj/typed/mongo/model/filters"
	"github.com/qianwj/typed/mongo/model/projections"
	"github.com/qianwj/typed/mongo/model/sorts"
	"github.com/qianwj/typed/mongo/operator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (p *Pipeline) AddFields(fields ...primitive.E) *Pipeline {
	p.append(operator.AddFields, primitive.D(fields))
	return p
}

func (p *Pipeline) Count(field string) *Pipeline {
	if field == "" {
		field = "count"
	}
	p.append(operator.Count, field)
	return p
}

func (p *Pipeline) GraphLookup(cond *lookup.GraphJoinCondition) *Pipeline {
	p.append(operator.GraphLookup, cond)
	return p
}

func (p *Pipeline) Group(id any, fields ...bson.Entry) *Pipeline {
	body := bson.D(bson.E("_id", id))
	for _, acc := range fields {
		body.Put(acc.Key, acc.Value)
	}
	p.append(operator.Group, body.Primitive())
	return p
}

func (p *Pipeline) Lookup(cond *lookup.JoinCondition) *Pipeline {
	p.append(operator.Lookup, cond)
	return p
}

func (p *Pipeline) Limit(limit int64) *Pipeline {
	p.append(operator.Limit, limit)
	return p
}

func (p *Pipeline) Match(filter *filters.Filter) *Pipeline {
	p.append(operator.Match, filter)
	return p
}

func (p *Pipeline) Project(projection *projections.Options) *Pipeline {
	p.append(operator.Project, projection)
	return p
}

func (p *Pipeline) Set(fields *bson.Map) *Pipeline {
	p.append(operator.Set, fields)
	return p
}

func (p *Pipeline) ShardedDataDistribution() *Pipeline {
	p.append(operator.ShardedDataDistribution, bson.M())
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

func (p *Pipeline) UnionWith(coll string, pipeline *Pipeline) *Pipeline {
	p.append(operator.UnionWith, bson.M(
		bson.E("coll", coll),
		bson.E("pipeline", pipeline.stages),
	))
	return p
}

func (p *Pipeline) Unset(fields ...string) *Pipeline {
	p.append(operator.Unset, fields)
	return p
}

func (p *Pipeline) Unwind(opts *UnwindOptions) *Pipeline {
	p.append(operator.Unwind, opts)
	return p
}
