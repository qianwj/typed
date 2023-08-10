package executor

import (
	"context"
	"github.com/qianwj/typed/mongo/model"
	"github.com/qianwj/typed/mongo/model/filters"
	"github.com/qianwj/typed/mongo/model/sorts"
	"github.com/qianwj/typed/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	raw "go.mongodb.org/mongo-driver/mongo"
	rawopts "go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type FindExecutor[D model.Document[I], I model.DocumentId] struct {
	readprefPrimary *raw.Collection
	readprefDefault *raw.Collection
	filter          *filters.Filter
	primary         bool
	opts            *rawopts.FindOptions
}

func NewFindExecutor[D model.Document[I], I model.DocumentId](
	readprefPrimary, readprefDefault *raw.Collection,
	filter *filters.Filter,
) *FindExecutor[D, I] {
	return &FindExecutor[D, I]{
		readprefPrimary: readprefPrimary,
		readprefDefault: readprefDefault,
		filter:          filter,
		opts:            rawopts.Find(),
	}
}

func (f *FindExecutor[D, I]) Primary() *FindExecutor[D, I] {
	f.primary = true
	return f
}

// AllowDiskUse sets the value for the AllowDiskUse field.
func (f *FindExecutor[D, I]) AllowDiskUse() *FindExecutor[D, I] {
	f.opts.SetAllowDiskUse(true)
	return f
}

// AllowPartialResults sets the value for the AllowPartialResults field.
func (f *FindExecutor[D, I]) AllowPartialResults() *FindExecutor[D, I] {
	f.opts.SetAllowPartialResults(true)
	return f
}

// BatchSize sets the value for the BatchSize field.
func (f *FindExecutor[D, I]) BatchSize(i int32) *FindExecutor[D, I] {
	f.opts.SetBatchSize(i)
	return f
}

// Collation sets the value for the Collation field.
func (f *FindExecutor[D, I]) Collation(collation *options.Collation) *FindExecutor[D, I] {
	f.opts.SetCollation((*rawopts.Collation)(collation))
	return f
}

// Comment sets the value for the Comment field.
func (f *FindExecutor[D, I]) Comment(comment string) *FindExecutor[D, I] {
	f.opts.SetComment(comment)
	return f
}

// CursorType sets the value for the CursorType field.
func (f *FindExecutor[D, I]) CursorType(ct options.CursorType) *FindExecutor[D, I] {
	f.opts.SetCursorType(rawopts.CursorType(ct))
	return f
}

// Hint sets the value for the Hint field.
func (f *FindExecutor[D, I]) Hint(index string) *FindExecutor[D, I] {
	f.opts.SetHint(index)
	return f
}

// Let sets the value for the Let field.
func (f *FindExecutor[D, I]) Let(let bson.M) *FindExecutor[D, I] {
	f.opts.SetLet(let)
	return f
}

func (f *FindExecutor[D, I]) Skip(skip int64) *FindExecutor[D, I] {
	f.opts.SetSkip(skip)
	return f
}

func (f *FindExecutor[D, I]) Limit(limit int64) *FindExecutor[D, I] {
	f.opts.SetLimit(limit)
	return f
}

// Max sets the value for the Max field.
func (f *FindExecutor[D, I]) Max(max int) *FindExecutor[D, I] {
	f.opts.SetMax(max)
	return f
}

// MaxAwaitTime sets the value for the MaxAwaitTime field.
func (f *FindExecutor[D, I]) MaxAwaitTime(d time.Duration) *FindExecutor[D, I] {
	f.opts.SetMaxAwaitTime(d)
	return f
}

// MaxTime specifies the max time to allow the query to run.
func (f *FindExecutor[D, I]) MaxTime(d time.Duration) *FindExecutor[D, I] {
	f.opts.SetMaxTime(d)
	return f
}

// Min sets the value for the Min field.
func (f *FindExecutor[D, I]) Min(min int) *FindExecutor[D, I] {
	f.opts.SetMin(min)
	return f
}

// NoCursorTimeout sets the value for the NoCursorTimeout field.
func (f *FindExecutor[D, I]) NoCursorTimeout() *FindExecutor[D, I] {
	f.opts.SetNoCursorTimeout(true)
	return f
}

// ReturnKey sets the value for the ReturnKey field.
func (f *FindExecutor[D, I]) ReturnKey() *FindExecutor[D, I] {
	f.opts.SetReturnKey(true)
	return f
}

// ShowRecordID sets the value for the ShowRecordID field.
func (f *FindExecutor[D, I]) ShowRecordID() *FindExecutor[D, I] {
	f.opts.SetShowRecordID(true)
	return f
}

func (f *FindExecutor[D, I]) Sort(sort *sorts.SortOptions) *FindExecutor[D, I] {
	f.opts.SetSort(sort.Marshal())
	return f
}

// OplogReplay sets the value for the OplogReplay field.
//
// Deprecated: This option has been deprecated in MongoDB version 4.4 and will be ignored by the server if it is set.
func (f *FindExecutor[D, I]) OplogReplay() *FindExecutor[D, I] {
	f.opts.SetOplogReplay(true)
	return f
}

// Snapshot sets the value for the Snapshot field.
//
// Deprecated: This option has been deprecated in MongoDB version 3.6 and removed in MongoDB version 4.0.
func (f *FindExecutor[D, I]) Snapshot() *FindExecutor[D, I] {
	f.opts.SetSnapshot(true)
	return f
}

func (f *FindExecutor[D, I]) Execute(ctx context.Context) ([]D, error) {
	var (
		data   []D
		err    error
		cursor *raw.Cursor
	)
	if f.primary {
		cursor, err = f.readprefPrimary.Find(ctx, f.filter.Marshal(), f.opts)
	} else {
		cursor, err = f.readprefDefault.Find(ctx, f.filter.Marshal(), f.opts)
	}
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (f *FindExecutor[D, I]) ExecuteTo(ctx context.Context, data any) error {
	var (
		err    error
		cursor *raw.Cursor
	)
	if f.primary {
		cursor, err = f.readprefPrimary.Find(ctx, f.filter.Marshal(), f.opts)
	} else {
		cursor, err = f.readprefDefault.Find(ctx, f.filter.Marshal(), f.opts)
	}
	if err != nil {
		return err
	}
	return cursor.All(ctx, &data)
}
