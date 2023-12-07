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
	"github.com/qianwj/typed/mongo/model/filters"
	"github.com/qianwj/typed/mongo/options"
	"go.mongodb.org/mongo-driver/mongo"
	rawopts "go.mongodb.org/mongo-driver/mongo/options"
)

type DeleteExecutor struct {
	coll   *mongo.Collection
	filter *filters.Filter
	opts   *rawopts.DeleteOptions
}

func newDeleteExecutor(coll *mongo.Collection, filter *filters.Filter) *DeleteExecutor {
	return &DeleteExecutor{
		coll:   coll,
		filter: filter,
		opts:   rawopts.Delete(),
	}
}

// Collation sets the value for the Collation field.
func (d *DeleteExecutor) Collation(c *options.Collation) *DeleteExecutor {
	d.opts.SetCollation(c.Raw())
	return d
}

// Comment sets the value for the Comment field.
func (d *DeleteExecutor) Comment(comment bson.UnorderedMap) *DeleteExecutor {
	d.opts.SetComment(comment)
	return d
}

// Hint sets the value for the Hint field.
func (d *DeleteExecutor) Hint(index string) *DeleteExecutor {
	d.opts.SetHint(index)
	return d
}

// Let sets the value for the Let field.
func (d *DeleteExecutor) Let(let bson.UnorderedMap) *DeleteExecutor {
	d.opts.SetLet(let)
	return d
}

func (d *DeleteExecutor) One(ctx context.Context) (int64, error) {
	res, err := d.coll.DeleteOne(ctx, d.filter, d.opts)
	if err != nil {
		return -1, err
	}
	return res.DeletedCount, nil
}

func (d *DeleteExecutor) Many(ctx context.Context) (int64, error) {
	res, err := d.coll.DeleteMany(ctx, d.filter, d.opts)
	if err != nil {
		return -1, err
	}
	return res.DeletedCount, nil
}
