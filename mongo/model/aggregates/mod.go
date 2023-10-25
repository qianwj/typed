package aggregates

import (
	"github.com/qianwj/typed/mongo/bson"
	"github.com/qianwj/typed/mongo/model/aggregates/group"
	"github.com/qianwj/typed/mongo/model/aggregates/lookup"
	"github.com/qianwj/typed/mongo/model/filters"
	"github.com/qianwj/typed/mongo/model/projections"
	"github.com/qianwj/typed/mongo/model/sorts"
)

const (
	Now = "$$NOW"
)

func Count(field string) *Pipeline {
	return New().Count(field)
}

func GraphLookup(cond *lookup.GraphJoinCondition) *Pipeline {
	return New().GraphLookup(cond)
}

func Group(id group.ID, fields ...group.Accumulator) *Pipeline {
	return New().Group(id, fields...)
}

func Limit(limit int64) *Pipeline {
	return New().Limit(limit)
}

func Lookup(cond *lookup.JoinCondition) *Pipeline {
	return New().Lookup(cond)
}

func Match(filter *filters.Filter) *Pipeline {
	return New().Match(filter)
}

func Project(projection *projections.Options) *Pipeline {
	return New().Project(projection)
}

func Set(fields *bson.Map) *Pipeline {
	return New().Set(fields)
}

func Skip(skip int64) *Pipeline {
	return New().Skip(skip)
}

func Sort(opts *sorts.Options) *Pipeline {
	return New().Sort(opts)
}

func SortByCount(expression any) *Pipeline {
	return New().SortByCount(expression)
}

func Unset(fields ...string) *Pipeline {
	return New().Unset(fields...)
}
