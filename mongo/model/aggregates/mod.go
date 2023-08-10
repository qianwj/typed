package aggregates

import (
	"github.com/qianwj/typed/mongo/model"
	"github.com/qianwj/typed/mongo/model/filters"
)

func Facet(facets ...*model.Facet) *Pipeline {
	return New().Facet(facets...)
}

func Match(filters ...*filters.Filter) *Pipeline {
	return New().Match(filters...)
}
