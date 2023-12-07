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
	"github.com/qianwj/typed/mongo/model/updates"
	"github.com/qianwj/typed/mongo/options"
	"go.mongodb.org/mongo-driver/mongo"
	rawopts "go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type FindOneAndUpdateExecutor[D model.Doc[I], I model.ID] struct {
	coll   *mongo.Collection
	filter *filters.Filter
	update *updates.Update
	opts   *rawopts.FindOneAndUpdateOptions
}

func newFindOneAndUpdateExecutor[D model.Doc[I], I model.ID](
	primary *mongo.Collection,
	filter *filters.Filter,
	update *updates.Update,
) *FindOneAndUpdateExecutor[D, I] {
	return &FindOneAndUpdateExecutor[D, I]{
		coll:   primary,
		filter: filter,
		update: update,
		opts:   rawopts.FindOneAndUpdate(),
	}
}

// BypassDocumentValidation sets the value for the BypassDocumentValidation field.
func (f *FindOneAndUpdateExecutor[D, I]) BypassDocumentValidation() *FindOneAndUpdateExecutor[D, I] {
	f.opts.SetBypassDocumentValidation(true)
	return f
}

// ArrayFilters sets the value for the ArrayFilters field.
func (f *FindOneAndUpdateExecutor[D, I]) ArrayFilters(af options.ArrayFilters) *FindOneAndUpdateExecutor[D, I] {
	f.opts.SetArrayFilters(af.Raw())
	return f
}

// Collation sets the value for the Collation field.
func (f *FindOneAndUpdateExecutor[D, I]) Collation(collation *options.Collation) *FindOneAndUpdateExecutor[D, I] {
	f.opts.SetCollation(collation.Raw())
	return f
}

// Comment sets the value for the Comment field.
func (f *FindOneAndUpdateExecutor[D, I]) Comment(comment interface{}) *FindOneAndUpdateExecutor[D, I] {
	f.opts.SetComment(comment)
	return f
}

// MaxTime sets the value for the MaxTime field.
//
// NOTE(benjirewis): MaxTime will be deprecated in a future release. The more general Timeout
// option may be used in its place to control the amount of time that a single operation can
// run before returning an error. MaxTime is ignored if Timeout is set on the client.
func (f *FindOneAndUpdateExecutor[D, I]) MaxTime(d time.Duration) *FindOneAndUpdateExecutor[D, I] {
	f.opts.SetMaxTime(d)
	return f
}

// Projection sets the value for the Projection field.
func (f *FindOneAndUpdateExecutor[D, I]) Projection(projection *projections.Options) *FindOneAndUpdateExecutor[D, I] {
	f.opts.SetProjection(projection)
	return f
}

func (f *FindOneAndUpdateExecutor[D, I]) ReturnBefore() *FindOneAndUpdateExecutor[D, I] {
	f.opts.SetReturnDocument(rawopts.Before)
	return f
}

func (f *FindOneAndUpdateExecutor[D, I]) ReturnAfter() *FindOneAndUpdateExecutor[D, I] {
	f.opts.SetReturnDocument(rawopts.After)
	return f
}

// Sort sets the value for the Sort field.
func (f *FindOneAndUpdateExecutor[D, I]) Sort(sort *sorts.Options) *FindOneAndUpdateExecutor[D, I] {
	f.opts.SetSort(sort)
	return f
}

// Upsert sets the value for the Upsert field.
func (f *FindOneAndUpdateExecutor[D, I]) Upsert() *FindOneAndUpdateExecutor[D, I] {
	f.opts.SetUpsert(true)
	return f
}

// Hint sets the value for the Hint field.
func (f *FindOneAndUpdateExecutor[D, I]) Hint(index string) *FindOneAndUpdateExecutor[D, I] {
	f.opts.SetHint(index)
	return f
}

// Let sets the value for the Let field.
func (f *FindOneAndUpdateExecutor[D, I]) Let(l bson.UnorderedMap) *FindOneAndUpdateExecutor[D, I] {
	f.opts.SetLet(l)
	return f
}

func (f *FindOneAndUpdateExecutor[D, I]) Execute(ctx context.Context) (D, error) {
	res := f.coll.FindOneAndUpdate(ctx, f.filter, f.update, f.opts)
	var data D
	if res.Err() != nil {
		return data, res.Err()
	}
	if err := res.Decode(&data); err != nil {
		return data, err
	}
	return data, nil
}

func (f *FindOneAndUpdateExecutor[D, I]) Collect(ctx context.Context, data any) error {
	res := f.coll.FindOneAndUpdate(ctx, f.filter, f.update, f.opts)
	if res.Err() != nil {
		return res.Err()
	}
	return res.Decode(data)
}
