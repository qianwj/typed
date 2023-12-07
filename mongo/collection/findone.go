// MIT License
//
// Copyright (c) 2022 qianwj
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package collection

import (
	"context"
	"github.com/qianwj/typed/mongo/model"
	"github.com/qianwj/typed/mongo/model/filters"
	"github.com/qianwj/typed/mongo/model/projections"
	"github.com/qianwj/typed/mongo/model/sorts"
	"github.com/qianwj/typed/mongo/options"
	raw "go.mongodb.org/mongo-driver/mongo"
	rawopts "go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type FindOneExecutor[D model.Doc[I], I model.ID] struct {
	readprefPrimary *raw.Collection
	readprefDefault *raw.Collection
	filter          *filters.Filter
	primary         bool
	opts            *rawopts.FindOneOptions
}

func newFindOneExecutor[D model.Doc[I], I model.ID](readprefPrimary, readprefDefault *raw.Collection, filter *filters.Filter) *FindOneExecutor[D, I] {
	return &FindOneExecutor[D, I]{
		readprefPrimary: readprefPrimary,
		readprefDefault: readprefDefault,
		filter:          filter,
		opts:            rawopts.FindOne(),
	}
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
	f.opts.SetCollation(collation.Raw())
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

func (f *FindOneExecutor[D, I]) Sort(sort *sorts.Options) *FindOneExecutor[D, I] {
	f.opts.SetSort(sort)
	return f
}

func (f *FindOneExecutor[D, I]) Projection(projection *projections.Options) *FindOneExecutor[D, I] {
	f.opts.SetProjection(projection)
	return f
}

// MaxAwaitTime sets the value for the MaxAwaitTime field.
//
// Deprecated: This option is not valid for a findOne operation, as no cursor is actually created.
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
		res = f.readprefPrimary.FindOne(ctx, f.filter, f.opts)
	} else {
		res = f.readprefDefault.FindOne(ctx, f.filter, f.opts)
	}
	if res.Err() != nil {
		return data, res.Err()
	}
	if err := res.Decode(&data); err != nil {
		return data, err
	}
	return data, nil
}

func (f *FindOneExecutor[D, I]) Collect(ctx context.Context, data any) error {
	var res *raw.SingleResult
	if f.primary {
		res = f.readprefPrimary.FindOne(ctx, f.filter, f.opts)
	} else {
		res = f.readprefDefault.FindOne(ctx, f.filter, f.opts)
	}
	if res.Err() != nil {
		return res.Err()
	}
	return res.Decode(data)
}
