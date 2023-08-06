package executor

import (
	"context"
	"github.com/qianwj/typed/mongo/model"
	"github.com/qianwj/typed/mongo/model/filter"
	"github.com/qianwj/typed/mongo/options"
	raw "go.mongodb.org/mongo-driver/mongo"
	rawopts "go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type FindOneExecutor[D model.Document[I], I model.DocumentId] struct {
	coll    *Collection[D, I]
	filter  *filter.Filter
	primary bool
	opts    *rawopts.FindOneOptions
}

func (f *FindOneExecutor[D, I]) Primary() *FindOneExecutor[D, I] {
	f.primary = true
	return f
}

// AllowPartialResults sets the value for the AllowPartialResults field.
func (f *FindOneExecutor[D, I]) AllowPartialResults() *FindOneExecutor[D, I] {
	f.opts.SetAllowPartialResults(true)
	return f
}

// Collation sets the value for the Collation field.
func (f *FindOneExecutor[D, I]) Collation(collation *options.Collation) *FindOneExecutor[D, I] {
	f.opts.SetCollation((*rawopts.Collation)(collation))
	return f
}

// Comment sets the value for the Comment field.
func (f *FindOneExecutor[D, I]) Comment(comment string) *FindOneExecutor[D, I] {
	f.opts.SetComment(comment)
	return f
}

// Hint sets the value for the Hint field.
func (f *FindOneExecutor[D, I]) Hint(index string) *FindOneExecutor[D, I] {
	f.opts.SetHint(index)
	return f
}

func (f *FindOneExecutor[D, I]) Skip(skip int64) *FindOneExecutor[D, I] {
	f.opts.SetSkip(skip)
	return f
}

// Max sets the value for the Max field.
func (f *FindOneExecutor[D, I]) Max(max int) *FindOneExecutor[D, I] {
	f.opts.SetMax(max)
	return f
}

// MaxTime specifies the max time to allow the query to run.
func (f *FindOneExecutor[D, I]) MaxTime(d time.Duration) *FindOneExecutor[D, I] {
	f.opts.SetMaxTime(d)
	return f
}

// Min sets the value for the Min field.
func (f *FindOneExecutor[D, I]) Min(min int) *FindOneExecutor[D, I] {
	f.opts.SetMin(min)
	return f
}

// ReturnKey sets the value for the ReturnKey field.
func (f *FindOneExecutor[D, I]) ReturnKey() *FindOneExecutor[D, I] {
	f.opts.SetReturnKey(true)
	return f
}

// ShowRecordID sets the value for the ShowRecordID field.
func (f *FindOneExecutor[D, I]) ShowRecordID() *FindOneExecutor[D, I] {
	f.opts.SetShowRecordID(true)
	return f
}

func (f *FindOneExecutor[D, I]) Sort(sort options.SortOptions) *FindOneExecutor[D, I] {
	f.opts.SetSort(sort)
	return f
}

// MaxAwaitTime sets the value for the MaxAwaitTime field.
func (f *FindOneExecutor[D, I]) MaxAwaitTime(d time.Duration) *FindOneExecutor[D, I] {
	f.opts.SetMaxAwaitTime(d)
	return f
}

// OplogReplay sets the value for the OplogReplay field.
//
// Deprecated: This option has been deprecated in MongoDB version 4.4 and will be ignored by the server if it is set.
func (f *FindOneExecutor[D, I]) OplogReplay() *FindOneExecutor[D, I] {
	f.opts.SetOplogReplay(true)
	return f
}

// Snapshot sets the value for the Snapshot field.
//
// Deprecated: This option has been deprecated in MongoDB version 3.6 and removed in MongoDB version 4.0.
func (f *FindOneExecutor[D, I]) Snapshot() *FindOneExecutor[D, I] {
	f.opts.SetSnapshot(true)
	return f
}

func (f *FindOneExecutor[D, I]) Execute(ctx context.Context) (D, error) {
	var (
		data D
		res  *raw.SingleResult
	)
	if f.primary {
		res = f.coll.primary.FindOne(ctx, f.filter.Marshal(), f.opts)
	} else {
		res = f.coll.defaultReadpref.FindOne(ctx, f.filter.Marshal(), f.opts)
	}
	if res.Err() != nil {
		return data, res.Err()
	}
	if err := res.Decode(&data); err != nil {
		return data, err
	}
	return data, nil
}

func (f *FindOneExecutor[D, I]) ExecuteTo(ctx context.Context, data any) error {
	var res *raw.SingleResult
	if f.primary {
		res = f.coll.primary.FindOne(ctx, f.filter.Marshal(), f.opts)
	} else {
		res = f.coll.defaultReadpref.FindOne(ctx, f.filter.Marshal(), f.opts)
	}
	if res.Err() != nil {
		return res.Err()
	}
	return res.Decode(&data)
}
