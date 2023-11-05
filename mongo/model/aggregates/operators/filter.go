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

package operators

import (
	"github.com/qianwj/typed/mongo/bson"
	"github.com/qianwj/typed/mongo/util"
)

type FilterSource struct {
	input bson.Array
	cond  any
	as    *string
	limit any
}

func NewFilterSource(input bson.Array, cond any) *FilterSource {
	return &FilterSource{
		input: input,
		cond:  cond,
	}
}

func (f *FilterSource) As(as string) *FilterSource {
	f.as = util.ToPtr(as)
	return f
}

func (f *FilterSource) Limit(limitExpr any) *FilterSource {
	f.limit = limitExpr
	return f
}

func (f *FilterSource) MarshalBSON() ([]byte, error) {
	m := bson.M(
		bson.E("input", f.input),
		bson.E("cond", f.cond),
	)
	if util.IsNonNil(f.as) {
		m["as"] = f.as
	}
	if util.IsNonNil(f.limit) {
		m["limit"] = f.limit
	}
	return m.Marshal()
}
