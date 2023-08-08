package options

import (
	"github.com/qianwj/typed/mongo/operator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	Desc SortOrder = -1
	Asc  SortOrder = 1
)

const (
	// NonTailable specifies that a cursor should close after retrieving the last data.
	NonTailable CursorType = iota
	// Tailable specifies that a cursor should not close when the last data is retrieved and can be resumed later.
	Tailable
	// TailableAwait specifies that a cursor should not close when the last data is retrieved and
	// that it should block for a certain amount of time for new data before returning no data.
	TailableAwait
)

type (
	CursorType  options.CursorType
	Collation   options.Collation
	SortOrder   int
	SortOptions bson.D
	SortOption  bson.E
	Projection  bson.M
	Page        struct {
		pageNo   int64
		pageSize int64
	}
)

func (opts SortOptions) Append(next SortOptions) SortOptions {
	return append(opts, next...)
}

func Ascending(field string) SortOptions {
	return SortOptions{bson.E{Key: field, Value: Asc}}
}

func Descending(field string) SortOptions {
	return SortOptions{bson.E{Key: field, Value: Desc}}
}

func Meta(field string) SortOptions {
	return SortOptions{
		{Key: field, Value: bson.M{operator.Meta: "textScore"}},
	}
}

func Need(field string) Projection {
	return Projection{
		field: 1,
	}
}

func Ignore(field string) Projection {
	return Projection{
		field: -1,
	}
}

func (p Projection) And(projection Projection) Projection {
	for k, v := range projection {
		p[k] = v
	}
	return p
}

func NewPage(pageNo, pageSize int64) *Page {
	return &Page{
		pageNo:   pageNo,
		pageSize: pageSize,
	}
}

func (p *Page) toOptions() *options.FindOptions {
	start := (p.pageNo - 1) * p.pageSize
	return options.Find().SetSkip(start).SetLimit(p.pageSize)
}
