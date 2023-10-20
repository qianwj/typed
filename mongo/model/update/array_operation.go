package update

import (
	"github.com/qianwj/typed/mongo/model/filters"
	"github.com/qianwj/typed/mongo/model/sorts"
	"github.com/qianwj/typed/mongo/operator"
	"go.mongodb.org/mongo-driver/bson"
)

func (u *Update) AddToSet(field string, value any) *Update {
	u.addToSet[field] = value
	return u
}

func (u *Update) AddEachToSet(field string, value []any) *Update {
	u.addToSet[field] = bson.M{operator.Each: value}
	return u
}

func (u *Update) Push(field string, value any) *Update {
	u.push[field] = value
	return u
}

func (u *Update) PushEach(field string, value []any) *Update {
	u.push[field] = bson.M{operator.Each: value}
	return u
}

func (u *Update) InsertAll(field string, value []any, index int) *Update {
	u.push[field] = bson.M{
		operator.Each:     value,
		operator.Position: index,
	}
	return u
}

func (u *Update) PushLimited(field string, value []string, limit int) *Update {
	u.push[field] = bson.M{
		operator.Each:  value,
		operator.Slice: limit,
	}
	return u
}

func (u *Update) PushSorted(field string, value []string, sort *sorts.SortOptions) *Update {
	u.push[field] = bson.M{
		operator.Each: value,
		operator.Sort: sort.Marshal(),
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

func (u *Update) PullConditioned(condition *filters.Filter) *Update {
	cond := condition.Raw()
	for _, e := range cond {
		u.pull[e.Key] = e.Value
	}
	return u
}

func (u *Update) PullAll(field string, value []any) *Update {
	u.pullAll[field] = value
	return u
}
