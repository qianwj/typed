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

type ReplaceExecutor[I any] struct {
	coll *mongo.Collection
	rep  any
	opts *rawoptions.ReplaceOptions
}

func newReplaceExecutor[I any](coll *mongo.Collection, rep any) *ReplaceExecutor[I] {
	return &ReplaceExecutor[I]{
		coll: coll,
		rep:  rep,
		opts: rawoptions.Replace(),
	}
}

// BypassDocumentValidation sets the value for the BypassDocumentValidation field.
func (r *ReplaceExecutor[I]) BypassDocumentValidation() *ReplaceExecutor[I] {
	r.opts.SetBypassDocumentValidation(true)
	return r
}

// Collation sets the value for the Collation field.
func (r *ReplaceExecutor[I]) Collation(c *options.Collation) *ReplaceExecutor[I] {
	r.opts.SetCollation(c.Raw())
	return r
}

// Comment sets the value for the Comment field.
func (r *ReplaceExecutor[I]) Comment(comment bson.UnorderedMap) *ReplaceExecutor[I] {
	r.opts.SetComment(comment)
	return r
}

// Hint sets the value for the Hint field.
func (r *ReplaceExecutor[I]) Hint(index string) *ReplaceExecutor[I] {
	r.opts.SetHint(index)
	return r
}

// Upsert sets the value for the Upsert field.
func (r *ReplaceExecutor[I]) Upsert() *ReplaceExecutor[I] {
	r.opts.SetUpsert(true)
	return r
}

// Let sets the value for the Let field.
func (r *ReplaceExecutor[I]) Let(l bson.UnorderedMap) *ReplaceExecutor[I] {
	r.opts.SetLet(l)
	return r
}

func (r *ReplaceExecutor[I]) Execute(ctx context.Context) (*updates.UpdateResult[I], error) {
	res, err := r.coll.ReplaceOne(ctx, r.rep, r.opts)
	return updates.FromUpdateResult[I](res), err
}
