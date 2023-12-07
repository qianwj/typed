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
	"github.com/qianwj/typed/mongo/model/filters"
	"go.mongodb.org/mongo-driver/mongo"
)

type Collection struct {
	coll *mongo.Collection
}

func newCollection(coll *mongo.Collection) *Collection {
	return &Collection{coll: coll}
}

func (c *Collection) Insert(docs ...any) *InsertExecutor {
	return newInsertExecutor(c.coll, docs...)
}

func (c *Collection) Delete(filter *filters.Filter) *DeleteExecutor {
	return newDeleteExecutor(c.coll, filter)
}

func (c *Collection) Replace(replacement any) *ReplaceExecutor[any] {
	return newReplaceExecutor[any](c.coll, replacement)
}
