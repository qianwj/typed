package update

import (
	"github.com/qianwj/typed/mongo/model/filter"
	operator2 "github.com/qianwj/typed/mongo/operator"
	"github.com/qianwj/typed/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
)

func (u *Update) AddToSet(field string, value any) *Update {
	u.addToSet[field] = value
	return u
}

func (u *Update) AddEachToSet(field string, value []any) *Update {
	u.addToSet[field] = bson.M{operator2.Each: value}
	return u
}

func (u *Update) Push(field string, value any) *Update {
	u.push[field] = value
	return u
}

func (u *Update) PushEach(field string, value []any) *Update {
	u.push[field] = bson.M{operator2.Each: value}
	return u
}

func (u *Update) InsertAll(field string, value []any, index int) *Update {
	u.push[field] = bson.M{
		operator2.Each:     value,
		operator2.Position: index,
	}
	return u
}

func (u *Update) PushLimited(field string, value []string, limit int) *Update {
	u.push[field] = bson.M{
		operator2.Each:  value,
		operator2.Slice: limit,
	}
	return u
}

func (u *Update) PushSorted(field string, value []string, sort options.SortOptions) *Update {
	u.push[field] = bson.M{
		operator2.Each: value,
		operator2.Sort: sort,
	}
	return u
}

func (u *Update) Pop(field string, pType PopType) *Update {
	u.pop[field] = pType
	return u
}

func (u *Update) Pull(field string, value any) *Update {
	u.pull[field] = value
	return u
}

func (u *Update) PullConditioned(condition *filter.Filter) *Update {
	cond := condition.Marshal()
	for _, e := range cond {
		u.pull[e.Key] = e.Value
	}
	return u
}

func (u *Update) PullAll(field string, value []any) *Update {
	u.pullAll[field] = value
	return u
}
