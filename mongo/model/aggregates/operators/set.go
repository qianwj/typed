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

// AllElementsTrue evaluates an array as a set and returns true if no element in the array is false. Otherwise, returns
// false. An empty array returns true.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/allElementsTrue/
func AllElementsTrue(expr any) bson.Entry {
	return bson.E(operator.AllElementsTrue, expr)
}

// AnyElementTrue evaluates an array as a set and returns true if any of the elements are true and false otherwise. An
// empty array returns false. array returns true.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/anyElementTrue/
func AnyElementTrue(expr any) bson.Entry {
	return bson.E(operator.AnyElementTrue, expr)
}

// SetDifference takes two sets and returns an array containing the elements that only exist in the first set; i.e.
// performs a `relative complement` of the second set relative to the first.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/setDifference/
func SetDifference(expr1, expr2 any) bson.Entry {
	return computeBoth(operator.SetDifference, expr1, expr2)
}

// SetEquals compares two or more arrays and returns true if they have the same distinct elements and false otherwise.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/setEquals/
func SetEquals(expr1, expr2 any) bson.Entry {
	return computeBoth(operator.SetEquals, expr1, expr2)
}

// SetIntersection takes two or more arrays and returns an array that contains the elements that appear in every input
// array.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/setIntersection/
func SetIntersection(exprs ...any) bson.Entry {
	return bson.E(operator.SetIntersection, bson.A(exprs...))
}

// SetIsSubset takes two arrays and returns true when the first array is a subset of the second, including when the
// first array equals the second array, and false otherwise.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/setIsSubset/
func SetIsSubset(expr1, expr2 any) bson.Entry {
	return computeBoth(operator.SetIsSubset, expr1, expr2)
}

// SetUnion takes two or more arrays and returns an array containing the elements that appear in any input array.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/setUnion/
func SetUnion(exprs ...any) bson.Entry {
	return bson.E(operator.SetUnion, bson.A(exprs...))
}
