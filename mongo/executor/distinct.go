package executor

import (
	"context"
	"github.com/qianwj/typed/mongo/options"
	raw "go.mongodb.org/mongo-driver/mongo"
	rawopts "go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type DistinctExecutor struct {
	readprefPrimary *raw.Collection
	readprefDefault *raw.Collection
	primary         bool
	field           string
	opts            *rawopts.DistinctOptions
}

func NewDistinctExecutor(readprefPrimary, readprefDefault *raw.Collection, field string) *DistinctExecutor {
	return &DistinctExecutor{
		readprefPrimary: readprefPrimary,
		readprefDefault: readprefDefault,
		field:           field,
		opts:            rawopts.Distinct(),
	}
}

func (do *DistinctExecutor) Primary() *DistinctExecutor {
	do.primary = true
	return do
}

// Collation sets the value for the Collation field.
func (do *DistinctExecutor) Collation(c *options.Collation) *DistinctExecutor {
	do.opts.SetCollation((*rawopts.Collation)(c))
	return do
}

// Comment sets the value for the Comment field.
func (do *DistinctExecutor) Comment(comment string) *DistinctExecutor {
	do.opts.SetComment(comment)
	return do
}

// MaxTime sets the value for the MaxTime field.
//
// NOTE(benjirewis): MaxTime will be deprecated in a future release. The more general Timeout
// option may be used in its place to control the amount of time that a single operation can
// run before returning an error. MaxTime is ignored if Timeout is set on the client.
func (do *DistinctExecutor) MaxTime(d time.Duration) *DistinctExecutor {
	do.opts.SetMaxTime(d)
	return do
}

func (do *DistinctExecutor) Execute(ctx context.Context) ([]any, error) {
	if do.primary {
		return do.readprefPrimary.Distinct(ctx, do.field, do.opts)
	}
	return do.readprefDefault.Distinct(ctx, do.field, do.opts)
}
