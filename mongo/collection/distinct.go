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
	"github.com/qianwj/typed/mongo/options"
	raw "go.mongodb.org/mongo-driver/mongo"
	rawopts "go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type DistinctExecutor struct {
	readprefPrimary *raw.Collection
	readprefDefault *raw.Collection
	primary         bool
	field           string
	opts            *rawopts.DistinctOptions
}

func newDistinctExecutor(readprefPrimary, readprefDefault *raw.Collection, field string) *DistinctExecutor {
	return &DistinctExecutor{
		readprefPrimary: readprefPrimary,
		readprefDefault: readprefDefault,
		field:           field,
		opts:            rawopts.Distinct(),
	}
}

func (do *DistinctExecutor) Primary() *DistinctExecutor {
	do.primary = true
	return do
}

// Collation sets the value for the Collation field.
func (do *DistinctExecutor) Collation(c *options.Collation) *DistinctExecutor {
	do.opts.SetCollation(c.Raw())
	return do
}

// Comment sets the value for the Comment field.
func (do *DistinctExecutor) Comment(comment bson.UnorderedMap) *DistinctExecutor {
	do.opts.SetComment(comment)
	return do
}

// MaxTime sets the value for the MaxTime field.
//
// NOTE(benjirewis): MaxTime will be deprecated in a future release. The more general Timeout
// option may be used in its place to control the amount of time that a single operation can
// run before returning an error. MaxTime is ignored if Timeout is set on the client.
func (do *DistinctExecutor) MaxTime(d time.Duration) *DistinctExecutor {
	do.opts.SetMaxTime(d)
	return do
}

func (do *DistinctExecutor) Execute(ctx context.Context) ([]any, error) {
	if do.primary {
		return do.readprefPrimary.Distinct(ctx, do.field, do.opts)
	}
	return do.readprefDefault.Distinct(ctx, do.field, do.opts)
}
