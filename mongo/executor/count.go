package executor

import (
	"context"
	"github.com/qianwj/typed/mongo/model"
	"github.com/qianwj/typed/mongo/model/filter"
	"github.com/qianwj/typed/mongo/options"
	rawopts "go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type CountExecutor[D model.Document[I], I model.DocumentId] struct {
	coll    *Collection[D, I]
	filter  *filter.Filter
	primary bool
	opts    *rawopts.CountOptions
}

func (c *CountExecutor[D, I]) Primary() *CountExecutor[D, I] {
	c.primary = true
	return c
}

// Collation sets the value for the Collation field.
func (c *CountExecutor[D, I]) Collation(collation *options.Collation) *CountExecutor[D, I] {
	c.opts.SetCollation((*rawopts.Collation)(collation))
	return c
}

// Comment sets the value for the Comment field.
func (c *CountExecutor[D, I]) Comment(comment string) *CountExecutor[D, I] {
	c.opts.SetComment(comment)
	return c
}

// Hint sets the value for the Hint field.
func (c *CountExecutor[D, I]) Hint(index string) *CountExecutor[D, I] {
	c.opts.SetHint(index)
	return c
}

// Limit sets the value for the Limit field.
func (c *CountExecutor[D, I]) Limit(i int64) *CountExecutor[D, I] {
	c.opts.SetLimit(i)
	return c
}

// MaxTime sets the value for the MaxTime field.
//
// NOTE(benjirewis): MaxTime will be deprecated in a future release. The more general Timeout
// option may be used in its place to control the amount of time that a single operation can
// run before returning an error. MaxTime is ignored if Timeout is set on the client.
func (c *CountExecutor[D, I]) MaxTime(d time.Duration) *CountExecutor[D, I] {
	c.opts.SetMaxTime(d)
	return c
}

// Skip sets the value for the Skip field.
func (c *CountExecutor[D, I]) Skip(i int64) *CountExecutor[D, I] {
	c.opts.SetSkip(i)
	return c
}

func (c *CountExecutor[D, I]) Execute(ctx context.Context) (int64, error) {
	if c.primary {
		return c.coll.primary.CountDocuments(ctx, c.filter, c.opts)
	} else {
		return c.coll.defaultReadpref.CountDocuments(ctx, c.filter, c.opts)
	}
}
