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

// ConcatArrays concatenates arrays to return the concatenated array.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/concatArrays/
func ConcatArrays(exprs ...any) bson.Entry {
	return bson.E(operator.ConcatArrays, bson.A(exprs...))
}

// Filter selects a subset of an array to return based on the specified condition. Returns an array with only those
// elements that match the condition. The returned elements are in the original order.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/filter/
func Filter(source *FilterSource) bson.Entry {
	return bson.E(operator.Filter, source)
}

// FirstN returns an aggregation of the first n elements within a group. The elements returned are meaningful only if
// in a specified sort order. If the group contains fewer than n elements, `$firstN` returns all elements in the group.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/firstN/
func FirstN(input, n any) bson.Entry {
	return bson.E(operator.FirstN, bson.D(
		bson.E("input", input),
		bson.E("n", n),
	))
}

// ArrayFirstN returns a specified number of elements from the beginning of an array.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/firstN-array-element/
func ArrayFirstN(n, input any) bson.Entry {
	return bson.E(operator.FirstN, bson.D(
		bson.E("n", n),
		bson.E("input", input),
	))
}

// In returns a boolean indicating whether a specified value is in an array.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/in/
func In(elemExpr, arrExpr any) bson.Entry {
	return computeBoth(operator.In, elemExpr, arrExpr)
}

// LastN returns an aggregation of the last n elements within a group. The elements returned are meaningful only if in
// a specified sort order. If the group contains fewer than n elements, `$lastN` returns all elements in the group.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/lastN/
func LastN(input, n any) bson.Entry {
	return bson.E(operator.LastN, bson.D(
		bson.E("input", input),
		bson.E("n", n),
	))
}

// ArrayLastN returns a specified number of elements from the end of an array.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/lastN-array-element/
func ArrayLastN(n, input any) bson.Entry {
	return bson.E(operator.LastN, bson.D(
		bson.E("n", n),
		bson.E("input", input),
	))
}

// Map applies an expression to each item in an array and returns an array with the applied results.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/map/
func Map(input, in any, as string) bson.Entry {
	return bson.E(operator.Map, bson.D(
		bson.E("input", input),
		bson.E("as", as),
		bson.E("in", in),
	))
}

// MaxN returns an aggregation of the maxmimum value n elements within a group. If the group contains fewer than n
// elements, $maxN returns all elements in the group.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/maxN/
func MaxN(input, n any) bson.Entry {
	return bson.E(operator.MaxN, bson.D(
		bson.E("input", input),
		bson.E("n", n),
	))
}

// ArrayMaxN returns the n largest values in an array.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/maxN-array-element/
func ArrayMaxN(n, input any) bson.Entry {
	return bson.E(operator.MaxN, bson.D(
		bson.E("n", n),
		bson.E("input", input),
	))
}

// MinN Returns an aggregation of the minimum value n elements within a group. If the group contains fewer than n
// elements, `$minN` returns all elements in the group.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/minN/
func MinN(input, n any) bson.Entry {
	return bson.E(operator.MinN, bson.D(
		bson.E("input", input),
		bson.E("n", n),
	))
}

// ArrayMinN returns the n smallest values in an array.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/minN-array-element/
func ArrayMinN(n, input any) bson.Entry {
	return bson.E(operator.MinN, bson.D(
		bson.E("n", n),
		bson.E("input", input),
	))
}

// ObjectToArray Converts a document to an array. The return array contains an element for each field/value pair in the
// original document. Each element in the return array is a document that contains two fields k and v:
//   - The k field contains the field name in the original document.
//   - The v field contains the value of the field in the original document.
//
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/objectToArray/
func ObjectToArray(expr any) bson.Entry {
	return bson.E(operator.ObjectToArray, expr)
}

// Reduce applies an expression to each element in an array and combines them into a single value.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/reduce/
func Reduce(input, initialVal, in any) bson.Entry {
	return bson.E(operator.Reduce, bson.D(
		bson.E("input", input),
		bson.E("initialValue", initialVal),
		bson.E("in", in),
	))
}
