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
	"github.com/qianwj/typed/mongo/operator"
)

// Accumulator defines a custom accumulator operator. Accumulators are operators that maintain their state (e.g. totals,
// maximums, minimums, and related data) as documents progress through the pipeline. Use the `$accumulator` operator to
// execute your own JavaScript functions to implement behavior not supported by the MongoDB Query Language.
// See also $function.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/accumulator/
func Accumulator(acc *AccumulatorOptions) bson.Entry {
	return bson.E(operator.Accumulator, acc)
}

// Function defines a custom aggregation function or expression in JavaScript.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/function/
func Function(code, lang string, args ...any) bson.Entry {
	return bson.E(operator.Function, bson.M(
		bson.E("code", code),
		bson.E("lang", lang),
		bson.E("args", bson.A(args...)),
	))
}
