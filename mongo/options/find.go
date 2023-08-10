package options

import (
	"go.mongodb.org/mongo-driver/mongo/options"
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
	CursorType options.CursorType
	Collation  options.Collation
	Page       struct {
		pageNo   int64
		pageSize int64
	}
)

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
