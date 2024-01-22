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
	"github.com/qianwj/typed/mongo/bson"
	"github.com/qianwj/typed/mongo/model"
	"github.com/qianwj/typed/mongo/model/filters"
	"github.com/qianwj/typed/mongo/model/projections"
	"github.com/qianwj/typed/mongo/model/sorts"
	"github.com/qianwj/typed/mongo/options"
	"go.mongodb.org/mongo-driver/mongo"
	rawopts "go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type FindOneAndDeleteExecutor[D model.Doc[I], I model.ID] struct {
	coll   *mongo.Collection
	filter *filters.Filter
	opts   *rawopts.FindOneAndDeleteOptions
}

func newFindOneAndDeleteExecutor[D model.Doc[I], I model.ID](
	primary *mongo.Collection,
	filter *filters.Filter,
) *FindOneAndDeleteExecutor[D, I] {
	return &FindOneAndDeleteExecutor[D, I]{
		coll:   primary,
		filter: filter,
		opts:   rawopts.FindOneAndDelete(),
	}
}

// Collation sets the value for the Collation field.
func (f *FindOneAndDeleteExecutor[D, I]) Collation(collation *options.Collation) *FindOneAndDeleteExecutor[D, I] {
	f.opts.SetCollation(collation.Raw())
	return f
}

// Comment sets the value for the Comment field.
func (f *FindOneAndDeleteExecutor[D, I]) Comment(comment interface{}) *FindOneAndDeleteExecutor[D, I] {
	f.opts.SetComment(comment)
	return f
}

// MaxTime sets the value for the MaxTime field.
//
// NOTE(benjirewis): MaxTime will be deprecated in a future release. The more general Timeout
// option may be used in its place to control the amount of time that a single operation can
// run before returning an error. MaxTime is ignored if Timeout is set on the client.
func (f *FindOneAndDeleteExecutor[D, I]) MaxTime(d time.Duration) *FindOneAndDeleteExecutor[D, I] {
	f.opts.SetMaxTime(d)
	return f
}

// Projection sets the value for the Projection field.
func (f *FindOneAndDeleteExecutor[D, I]) Projection(projection *projections.Options) *FindOneAndDeleteExecutor[D, I] {
	f.opts.SetProjection(projection)
	return f
}

// Sort sets the value for the Sort field.
func (f *FindOneAndDeleteExecutor[D, I]) Sort(sort *sorts.Options) *FindOneAndDeleteExecutor[D, I] {
	f.opts.SetSort(sort)
	return f
}

// Hint sets the value for the Hint field.
func (f *FindOneAndDeleteExecutor[D, I]) Hint(index string) *FindOneAndDeleteExecutor[D, I] {
	f.opts.SetHint(index)
	return f
}

// Let sets the value for the Let field.
func (f *FindOneAndDeleteExecutor[D, I]) Let(l bson.UnorderedMap) *FindOneAndDeleteExecutor[D, I] {
	f.opts.SetLet(l)
	return f
}

func (f *FindOneAndDeleteExecutor[D, I]) Execute(ctx context.Context) (D, error) {
	res := f.coll.FindOneAndDelete(ctx, f.filter, f.opts)
	var data D
	if res.Err() != nil {
		return data, res.Err()
	}
	if err := res.Decode(&data); err != nil {
		return data, err
	}
	return data, nil
}

func (f *FindOneAndDeleteExecutor[D, I]) Collect(ctx context.Context, data any) error {
	res := f.coll.FindOneAndDelete(ctx, f.filter, f.opts)
	if res.Err() != nil {
		return res.Err()
	}
	return res.Decode(data)
}