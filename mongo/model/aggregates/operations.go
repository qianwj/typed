package aggregates

import (
	"github.com/qianwj/typed/mongo/model/filters"
	"github.com/qianwj/typed/mongo/operator"
	"github.com/qianwj/typed/mongo/util"
	"go.mongodb.org/mongo-driver/bson"
)

func (p *Pipeline) Match(query ...*filters.Filter) *Pipeline {
	p.put(operator.Match, util.Map(query, func(f *filters.Filter) bson.D {
		return f.Marshal()
	}))
	return p
}
