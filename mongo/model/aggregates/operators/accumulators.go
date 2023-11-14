/** MIT License
 *
 * Copyright (c) 2022-2024 qianwj
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package operators

import (
	"github.com/qianwj/typed/mongo/bson"
	"github.com/qianwj/typed/mongo/model/aggregates/expressions"
	"github.com/qianwj/typed/mongo/model/sorts"
	"github.com/qianwj/typed/mongo/operator"
)

// AddToSet returns an array of all unique values that results from applying an expression to each document in a group.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/addToSet/
func AddToSet(expression any) bson.Entry {
	return bson.E(operator.AddToSet, expression)
}

// Avg returns the average value of the numeric values. `$avg` ignores non-numeric values.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/avg/
func Avg(field string, exprs ...any) bson.Entry {
	return accMulti(field, operator.Avg, exprs...)
}

// Bottom returns the bottom element within a group according to the specified sort order.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/bottom/
func Bottom(sortBy *sorts.Options, output []string) bson.Entry {
	return bson.E(operator.Bottom, bson.M(
		bson.E("sortBy", sortBy),
		bson.E("output", output),
	))
}

// BottomN returns an aggregation of the bottom n elements within a group, according to the specified sort order. If
// the group contains fewer than n elements, $bottomN returns all elements in the group.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/bottomN/
func BottomN[I expressions.Integer](sortBy *sorts.Options, output []string, n I) bson.Entry {
	return bson.E(operator.BottomN, bson.M(
		bson.E("sortBy", sortBy),
		bson.E("output", output),
		bson.E("n", n),
	))
}

// Count returns the number of documents in a group.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/count-accumulator/
func Count() bson.Entry {
	return bson.E(operator.Count, bson.M())
}

// First returns the result of an expression for the first document in a group of documents. Only meaningful when
// documents are in a defined order.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/first/
func First[E expressions.Expression](expr E) bson.Entry {
	return bson.E(operator.First, expr)
}

// Last returns the result of an expression for the last document in a group of documents. Only meaningful when
// documents are in a defined order.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/last/
func Last[E expressions.Expression](expr E) bson.Entry {
	return bson.E(operator.Last, expr)
}

// Max returns the maximum value. `$max` compares both value and type, using the specified BSON comparison order for
// values of different types.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/max/
func Max(field string, exprs ...any) bson.Entry {
	return accMulti(field, operator.Max, exprs...)
}

// Median returns an approximation of the median, the 50th percentile, as a scalar value.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/median/
func Median(method string, inputs ...any) bson.Entry {
	val := bson.M(bson.E("method", method))
	if len(inputs) == 1 {
		val["input"] = inputs[0]
	} else {
		val["input"] = bson.A(inputs...)
	}
	return bson.E(operator.Median, val)
}

// Min returns the minimum value. `$min` compares both value and type, using the specified BSON comparison order for
// values of different types.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/min/
func Min(field string, exprs ...any) bson.Entry {
	return accMulti(field, operator.Min, exprs...)
}

// Push returns an array of all values that result from applying an expression to documents.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/push/
func Push(expr any) bson.Entry {
	return bson.E(operator.Push, expr)
}

// Sum calculates and returns the collective sum of numeric values. `$sum` ignores non-numeric values.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/sum/
func Sum(field string, exprs ...any) bson.Entry {
	return accMulti(field, operator.Sum, exprs...)
}

// Top returns the top element within a group according to the specified sort order.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/top/
func Top(sortBy *sorts.Options, output []string) bson.Entry {
	return bson.E(operator.Top, bson.M(
		bson.E("sortBy", sortBy),
		bson.E("output", output),
	))
}

// TopN returns an aggregation of the top n elements within a group, according to the specified sort order. If the
// group contains fewer than n elements, $topN returns all elements in the group.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/topN/
func TopN[I expressions.Integer](sortBy *sorts.Options, output []string, n I) bson.Entry {
	return bson.E(operator.TopN, bson.M(
		bson.E("sortBy", sortBy),
		bson.E("output", output),
		bson.E("n", n),
	))
}
