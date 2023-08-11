package aggregates

import (
	"github.com/qianwj/typed/mongo/model/filters"
	"github.com/qianwj/typed/mongo/operator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Matcher struct {
	f *filters.Filter
}

func Match(filter *filters.Filter) *Matcher {
	return &Matcher{f: filter}
}

func (m *Matcher) Marshal() primitive.D {
	return primitive.D{
		{Key: operator.Match, Value: m.f.Marshal()},
	}
}

func (m *Matcher) ToMap() primitive.M {
	return primitive.M{
		operator.Match: m.f.ToMap(),
	}
}
