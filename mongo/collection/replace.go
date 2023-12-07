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
	"github.com/qianwj/typed/mongo/model/updates"
	"github.com/qianwj/typed/mongo/options"
	"go.mongodb.org/mongo-driver/mongo"
	rawoptions "go.mongodb.org/mongo-driver/mongo/options"
)

type ReplaceOneExecutor[I any] struct {
	coll *mongo.Collection
	rep  any
	opts *rawoptions.ReplaceOptions
}

func newReplaceOneExecutor[I any](coll *mongo.Collection, rep any) *ReplaceOneExecutor[I] {
	return &ReplaceOneExecutor[I]{
		coll: coll,
		rep:  rep,
		opts: rawoptions.Replace(),
	}
}

// BypassDocumentValidation sets the value for the BypassDocumentValidation field.
func (r *ReplaceOneExecutor[I]) BypassDocumentValidation() *ReplaceOneExecutor[I] {
	r.opts.SetBypassDocumentValidation(true)
	return r
}

// Collation sets the value for the Collation field.
func (r *ReplaceOneExecutor[I]) Collation(c *options.Collation) *ReplaceOneExecutor[I] {
	r.opts.SetCollation(c.Raw())
	return r
}

// Comment sets the value for the Comment field.
func (r *ReplaceOneExecutor[I]) Comment(comment bson.UnorderedMap) *ReplaceOneExecutor[I] {
	r.opts.SetComment(comment)
	return r
}

// Hint sets the value for the Hint field.
func (r *ReplaceOneExecutor[I]) Hint(index string) *ReplaceOneExecutor[I] {
	r.opts.SetHint(index)
	return r
}

// Upsert sets the value for the Upsert field.
func (r *ReplaceOneExecutor[I]) Upsert() *ReplaceOneExecutor[I] {
	r.opts.SetUpsert(true)
	return r
}

// Let sets the value for the Let field.
func (r *ReplaceOneExecutor[I]) Let(l bson.UnorderedMap) *ReplaceOneExecutor[I] {
	r.opts.SetLet(l)
	return r
}

func (r *ReplaceOneExecutor[I]) Execute(ctx context.Context) (*updates.UpdateResult[I], error) {
	res, err := r.coll.ReplaceOne(ctx, r.rep, r.opts)
	return updates.FromUpdateResult[I](res), err
}
