package executor

import (
	"context"
	"github.com/qianwj/typed/mongo/bson"
	"github.com/qianwj/typed/mongo/model/aggregates"
	"github.com/qianwj/typed/mongo/options"
	rawbson "go.mongodb.org/mongo-driver/bson"
	raw "go.mongodb.org/mongo-driver/mongo"
	rawopts "go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type AggregateExecutor[D bson.Doc[I], I bson.ID] struct {
	readprefPrimary *raw.Collection
	readprefDefault *raw.Collection
	pipe            *aggregates.Pipeline
	primary         bool
	opts            *rawopts.AggregateOptions
}

func NewAggregateExecutor[D bson.Doc[I], I bson.ID](
	readprefPrimary, readprefDefault *raw.Collection,
	pipe *aggregates.Pipeline,
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
func (a *AggregateExecutor[D, I]) Let(let rawbson.M) *AggregateExecutor[D, I] {
	a.opts.SetLet(let)
	return a
}

// Custom sets the value for the Custom field. Key-value pairs of the BSON map should correlate
// with desired option names and values. Values must be Marshalable. Custom options may conflict
// with non-custom options, and custom options bypass client-side validation. Prefer using non-custom
// options where possible.
func (a *AggregateExecutor[D, I]) Custom(c rawbson.M) *AggregateExecutor[D, I] {
	a.opts.SetCustom(c)
	return a
}

func (a *AggregateExecutor[D, I]) Collect(ctx context.Context, result any) error {
	var (
		err    error
		cursor *raw.Cursor
	)
	if a.primary {
		cursor, err = a.readprefPrimary.Aggregate(ctx, a.pipe.Stages(), a.opts)
	} else {
		cursor, err = a.readprefDefault.Aggregate(ctx, a.pipe.Stages(), a.opts)
	}
	if err != nil {
		return err
	}
	return cursor.All(ctx, result)
}

type DatabaseAggregateExecutor struct {
	readprefPrimary *raw.Database
	readprefDefault *raw.Database
	pipe            *aggregates.Pipeline
	primary         bool
	opts            *rawopts.AggregateOptions
}

func NewDatabaseAggregateExecutor(readprefPrimary, readprefDefault *raw.Database, pipe *aggregates.Pipeline) *DatabaseAggregateExecutor {
	return &DatabaseAggregateExecutor{
		readprefPrimary: readprefPrimary,
		readprefDefault: readprefDefault,
		pipe:            pipe,
		opts:            rawopts.Aggregate(),
	}
}

func (a *DatabaseAggregateExecutor) Primary() *DatabaseAggregateExecutor {
	a.primary = true
	return a
}

// AllowDiskUse sets the value for the AllowDiskUse field.
func (a *DatabaseAggregateExecutor) AllowDiskUse() *DatabaseAggregateExecutor {
	a.opts.SetAllowDiskUse(true)
	return a
}

// BatchSize sets the value for the BatchSize field.
func (a *DatabaseAggregateExecutor) BatchSize(i int32) *DatabaseAggregateExecutor {
	a.opts.SetBatchSize(i)
	return a
}

// BypassDocumentValidation sets the value for the BypassDocumentValidation field.
func (a *DatabaseAggregateExecutor) BypassDocumentValidation() *DatabaseAggregateExecutor {
	a.opts.SetBypassDocumentValidation(true)
	return a
}

// Collation sets the value for the Collation field.
func (a *DatabaseAggregateExecutor) Collation(c *options.Collation) *DatabaseAggregateExecutor {
	a.opts.SetCollation((*rawopts.Collation)(c))
	return a
}

// MaxTime sets the value for the MaxTime field.
//
// NOTE(benjirewis): MaxTime will be deprecated in a future release. The more general Timeout
// option may be used in its place to control the amount of time that a single operation can
// run before returning an error. MaxTime is ignored if Timeout is set on the client.
func (a *DatabaseAggregateExecutor) MaxTime(d time.Duration) *DatabaseAggregateExecutor {
	a.opts.SetMaxTime(d)
	return a
}

// MaxAwaitTime sets the value for the MaxAwaitTime field.
func (a *DatabaseAggregateExecutor) MaxAwaitTime(d time.Duration) *DatabaseAggregateExecutor {
	a.opts.SetMaxAwaitTime(d)
	return a
}

// Comment sets the value for the Comment field.
func (a *DatabaseAggregateExecutor) Comment(s string) *DatabaseAggregateExecutor {
	a.opts.SetComment(s)
	return a
}

// Hint sets the value for the Hint field.
func (a *DatabaseAggregateExecutor) Hint(index string) *DatabaseAggregateExecutor {
	a.opts.SetHint(index)
	return a
}

// Let sets the value for the Let field.
func (a *DatabaseAggregateExecutor) Let(let rawbson.M) *DatabaseAggregateExecutor {
	a.opts.SetLet(let)
	return a
}

// Custom sets the value for the Custom field. Key-value pairs of the BSON map should correlate
// with desired option names and values. Values must be Marshalable. Custom options may conflict
// with non-custom options, and custom options bypass client-side validation. Prefer using non-custom
// options where possible.
func (a *DatabaseAggregateExecutor) Custom(c rawbson.M) *DatabaseAggregateExecutor {
	a.opts.SetCustom(c)
	return a
}

func (a *DatabaseAggregateExecutor) Collect(ctx context.Context, result any) error {
	var (
		err    error
		cursor *raw.Cursor
	)
	if a.primary {
		cursor, err = a.readprefPrimary.Aggregate(ctx, a.pipe.Stages(), a.opts)
	} else {
		cursor, err = a.readprefDefault.Aggregate(ctx, a.pipe.Stages(), a.opts)
	}
	if err != nil {
		return err
	}
	return cursor.All(ctx, result)
}
