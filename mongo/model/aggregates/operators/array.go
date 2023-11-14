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
	"github.com/qianwj/typed/mongo/model/aggregates/expressions"
	"github.com/qianwj/typed/mongo/model/sorts"
	"github.com/qianwj/typed/mongo/operator"
)

// ArrayElemAt returns the element at the specified array index.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/arrayElemAt/
func ArrayElemAt[A expressions.Array, I expressions.Integer](expr A, idx I) bson.Entry {
	return bson.E(operator.ArrayElemAt, bson.A(expr, idx))
}

// ArrayToObject converts an array into a single document; the array must be either:
//   - An array of two-element arrays where the first element is the field name, and the second element is the field value:
//   - An array of documents that contains two fields, k and v where: The k field contains the field name, The v field
//     contains the value of the field.
//
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/arrayToObject/
func ArrayToObject[A expressions.Array](expr A) bson.Entry {
	return bson.E(operator.ArrayToObject, expr)
}

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
func FirstN[A expressions.Array, I expressions.Integer](input A, n I) bson.Entry {
	return bson.E(operator.FirstN, bson.D(
		bson.E("input", input),
		bson.E("n", n),
	))
}

// ArrayFirstN returns a specified number of elements from the beginning of an array.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/firstN-array-element/
func ArrayFirstN[A expressions.Array, I expressions.Integer](n I, input A) bson.Entry {
	return bson.E(operator.FirstN, bson.D(
		bson.E("n", n),
		bson.E("input", input),
	))
}

// In returns a boolean indicating whether a specified value is in an array.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/in/
func In[E expressions.Expression, A expressions.Array](elem E, arr A) bson.Entry {
	return computeBoth(operator.In, elem, arr)
}

// IndexOfArray searches an array for an occurrence of a specified value and returns the array index of the first
// occurrence. Array indexes start at zero.
// examples :
// ```IndexOfArray(bson.A("1", "2", "3", "4"), "1") // search all```
// ```IndexOfArray(bson.A("1", "2", "3", "4"), "1", 1) // search from ["2", "3"]```
// ```IndexOfArray(bson.A("1", "2", "3", "4"), "1", 1, 3) // search from ["2", "3", "4"]```
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/indexOfArray/
func IndexOfArray[A expressions.Array, E expressions.Expression, I bson.Int](expr A, search E, ranges ...I) bson.Entry {
	val := bson.A(expr, search)
	for _, n := range ranges {
		val = val.Append(n)
	}
	return bson.E(operator.IndexOfArray, val)
}

// IsArray determines if the operand is an array. Returns a boolean.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/isArray/
func IsArray[A expressions.Array](expr A) bson.Entry {
	return bson.E(operator.IsArray, expr)
}

// LastN returns an aggregation of the last n elements within a group. The elements returned are meaningful only if in
// a specified sort order. If the group contains fewer than n elements, `$lastN` returns all elements in the group.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/lastN/
func LastN[A expressions.Array, I expressions.Integer](input A, n I) bson.Entry {
	return bson.E(operator.LastN, bson.D(
		bson.E("input", input),
		bson.E("n", n),
	))
}

// ArrayLastN returns a specified number of elements from the end of an array.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/lastN-array-element/
func ArrayLastN[A expressions.Array, I expressions.Integer](n I, input A) bson.Entry {
	return bson.E(operator.LastN, bson.D(
		bson.E("n", n),
		bson.E("input", input),
	))
}

// OrderedMap applies an expression to each item in an array and returns an array with the applied results.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/map/
func OrderedMap[A expressions.Array, E expressions.Expression](input A, in E, as string) bson.Entry {
	return bson.E(operator.OrderedMap, bson.D(
		bson.E("input", input),
		bson.E("as", as),
		bson.E("in", in),
	))
}

// MaxN returns an aggregation of the maxmimum value n elements within a group. If the group contains fewer than n
// elements, $maxN returns all elements in the group.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/maxN/
func MaxN[A expressions.Array, I expressions.Integer](input A, n I) bson.Entry {
	return bson.E(operator.MaxN, bson.D(
		bson.E("input", input),
		bson.E("n", n),
	))
}

// ArrayMaxN returns the n largest values in an array.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/maxN-array-element/
func ArrayMaxN[A expressions.Array, I expressions.Integer](n I, input A) bson.Entry {
	return bson.E(operator.MaxN, bson.D(
		bson.E("n", n),
		bson.E("input", input),
	))
}

// MinN Returns an aggregation of the minimum value n elements within a group. If the group contains fewer than n
// elements, `$minN` returns all elements in the group.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/minN/
func MinN[A expressions.Array, I expressions.Integer](input A, n I) bson.Entry {
	return bson.E(operator.MinN, bson.D(
		bson.E("input", input),
		bson.E("n", n),
	))
}

// ArrayMinN returns the n smallest values in an array.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/minN-array-element/
func ArrayMinN[A expressions.Array, I expressions.Integer](n I, input A) bson.Entry {
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
func ObjectToArray[E expressions.Expression](expr E) bson.Entry {
	return bson.E(operator.ObjectToArray, expr)
}

// Range returns an array whose elements are a generated sequence of numbers. `$range` generates the sequence from the
// specified starting number by successively incrementing the starting number by the specified step value up to but not
// including the end point.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/range/
func Range[I expressions.Integer](start, end I, nonZeroStep ...I) bson.Entry {
	val := bson.A(start, end)
	if len(nonZeroStep) > 0 {
		val = val.Append(nonZeroStep[0])
	}
	return bson.E(operator.Range, val)
}

// Reduce applies an expression to each element in an array and combines them into a single value.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/reduce/
func Reduce[A expressions.Array, E expressions.Expression](input A, initialVal E, in E) bson.Entry {
	return bson.E(operator.Reduce, bson.D(
		bson.E("input", input),
		bson.E("initialValue", initialVal),
		bson.E("in", in),
	))
}

// ReverseArray accepts an array expression as an argument and returns an array with the elements in reverse order.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/reverseArray/
func ReverseArray[A expressions.Array](expr A) bson.Entry {
	return bson.E(operator.ReverseArray, expr)
}

// Size counts and returns the total number of items in an array.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/size/
func Size[A expressions.Array](expr A) bson.Entry {
	return bson.E(operator.Size, expr)
}

// Slice returns a subset of an array.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/slice/
func Slice[A expressions.Array, I expressions.Integer](expr A, n I, position ...I) bson.Entry {
	if len(position) > 0 {
		return bson.E(operator.Slice, bson.A(expr, position[0], n))
	}
	return bson.E(operator.Slice, bson.A(expr, n))
}

// SortArray sorts an array based on its elements. The sort order is user specified.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/sortArray/
func SortArray[A expressions.Array](input A, sortBy *sorts.Options) bson.Entry {
	return bson.E(operator.SortArray, bson.M(
		bson.E("input", input),
		bson.E("sortBy", sortBy),
	))
}

// Zip transposes an array of input arrays so that the first element of the output array would be an array containing,
// the first element of the first input array, the first element of the second input array, etc.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/zip/
func Zip[A expressions.Array](inputs []any, useLongestLength bool, defaults A) bson.Entry {
	return bson.E(operator.Zip, bson.M(
		bson.E("inputs", bson.A(inputs...)),
		bson.E("useLongestLength", useLongestLength),
		bson.E("defaults", defaults),
	))
}
