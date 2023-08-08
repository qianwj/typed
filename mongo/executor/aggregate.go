package executor

import (
	"context"
	"github.com/qianwj/typed/mongo/model"
	"github.com/qianwj/typed/mongo/model/aggregate"
	"github.com/qianwj/typed/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	raw "go.mongodb.org/mongo-driver/mongo"
	rawopts "go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type AggregateExecutor[D model.Document[I], I model.DocumentId] struct {
	readprefPrimary *raw.Collection
	readprefDefault *raw.Collection
	pipe            aggregate.Pipeline
	primary         bool
	opts            *rawopts.AggregateOptions
}

func NewAggregateExecutor[D model.Document[I], I model.DocumentId](
	readprefPrimary, readprefDefault *raw.Collection,
	pipe aggregate.Pipeline,
) *AggregateExecutor[D, I] {
	return &AggregateExecutor[D, I]{
		readprefPrimary: readprefPrimary,
		readprefDefault: readprefDefault,
		pipe:            pipe,
		opts:            rawopts.Aggregate(),
	}
}

func (a *AggregateExecutor[D, I]) Primary() *AggregateExecutor[D, I] {
	a.primary = true
	return a
}

// AllowDiskUse sets the value for the AllowDiskUse field.
func (a *AggregateExecutor[D, I]) AllowDiskUse() *AggregateExecutor[D, I] {
	a.opts.SetAllowDiskUse(true)
	return a
}

// BatchSize sets the value for the BatchSize field.
func (a *AggregateExecutor[D, I]) BatchSize(i int32) *AggregateExecutor[D, I] {
	a.opts.SetBatchSize(i)
	return a
}

// BypassDocumentValidation sets the value for the BypassDocumentValidation field.
func (a *AggregateExecutor[D, I]) BypassDocumentValidation() *AggregateExecutor[D, I] {
	a.opts.SetBypassDocumentValidation(true)
	return a
}

// Collation sets the value for the Collation field.
func (a *AggregateExecutor[D, I]) Collation(c *options.Collation) *AggregateExecutor[D, I] {
	a.opts.SetCollation((*rawopts.Collation)(c))
	return a
}

// MaxTime sets the value for the MaxTime field.
//
// NOTE(benjirewis): MaxTime will be deprecated in a future release. The more general Timeout
// option may be used in its place to control the amount of time that a single operation can
// run before returning an error. MaxTime is ignored if Timeout is set on the client.
func (a *AggregateExecutor[D, I]) MaxTime(d time.Duration) *AggregateExecutor[D, I] {
	a.opts.SetMaxTime(d)
	return a
}

// MaxAwaitTime sets the value for the MaxAwaitTime field.
func (a *AggregateExecutor[D, I]) MaxAwaitTime(d time.Duration) *AggregateExecutor[D, I] {
	a.opts.SetMaxAwaitTime(d)
	return a
}

// Comment sets the value for the Comment field.
func (a *AggregateExecutor[D, I]) Comment(s string) *AggregateExecutor[D, I] {
	a.opts.SetComment(s)
	return a
}

// Hint sets the value for the Hint field.
func (a *AggregateExecutor[D, I]) Hint(index string) *AggregateExecutor[D, I] {
	a.opts.SetHint(index)
	return a
}

// Let sets the value for the Let field.
func (a *AggregateExecutor[D, I]) Let(let bson.M) *AggregateExecutor[D, I] {
	a.opts.SetLet(let)
	return a
}

// Custom sets the value for the Custom field. Key-value pairs of the BSON map should correlate
// with desired option names and values. Values must be Marshalable. Custom options may conflict
// with non-custom options, and custom options bypass client-side validation. Prefer using non-custom
// options where possible.
func (a *AggregateExecutor[D, I]) Custom(c bson.M) *AggregateExecutor[D, I] {
	a.opts.SetCustom(c)
	return a
}

func (a *AggregateExecutor[D, I]) ExecuteTo(ctx context.Context, result interface{}) error {
	var (
		err    error
		cursor *raw.Cursor
	)
	if a.primary {
		cursor, err = a.readprefPrimary.Aggregate(ctx, a.pipe.Marshal(), a.opts)
	} else {
		cursor, err = a.readprefDefault.Aggregate(ctx, a.pipe.Marshal(), a.opts)
	}
	if err != nil {
		return err
	}
	return cursor.All(ctx, &result)
}
