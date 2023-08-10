package aggregates

import "github.com/qianwj/typed/mongo/model/filters"

func Match(filters ...*filters.Filter) *Pipeline {
	return New().Match(filters...)
}
