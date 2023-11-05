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

package accumulators

import (
	"github.com/qianwj/typed/mongo/bson"
	"github.com/qianwj/typed/mongo/operator"
)

// Abs returns the absolute value of a number.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/abs/
func Abs(numOrExpr any) bson.Entry {
	return bson.E(operator.Abs, numOrExpr)
}

// Add adds numbers together or adds numbers and a date. If one of the arguments is a date, `$add` treats the other
// arguments as milliseconds to add to the date.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/add/
func Add(exprs ...any) bson.Entry {
	return bson.E(operator.Add, bson.A(exprs...))
}

// AddToSet returns an array of all unique values that results from applying an expression to each document in a group.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/addToSet/
func AddToSet(expression any) bson.Entry {
	return bson.E(operator.AddToSet, expression)
}

// And Evaluates one or more expressions and returns true if all of the expressions are true or if run with no argument
// expressions. Otherwise, `$and` returns false.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/and/
func And(exprs ...any) bson.Entry {
	return bson.E(operator.AddToSet, bson.A(exprs...))
}

// Avg returns the average value of the numeric values. `$avg` ignores non-numeric values.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/avg/
func Avg(expression any) bson.Entry {
	return bson.E(operator.Avg, expression)
}

// Ceil returns the smallest integer greater than or equal to the specified number.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/ceil/
func Ceil(numOrExpr any) bson.Entry {
	return bson.E(operator.Ceil, numOrExpr)
}

// Cmp Compares two values and returns:
//   - -1 if the first value is less than the second.
//   - 1 if the first value is greater than the second.
//   - 0 if the two values are equivalent.
//
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/cmp/
func Cmp(expr1, expr2 any) bson.Entry {
	return bson.E(operator.Cmp, bson.A(expr1, expr2))
}

// Convert converts a value to a specified type.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/convert/
func Convert(c *Converter) bson.Entry {
	return bson.E(operator.Convert, c)
}

// Count returns the number of documents in a group.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/count-accumulator/
func Count() bson.Entry {
	return bson.E(operator.Count, bson.M())
}

// DateAdd increments a `Date()` object by a specified number of time units.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/dateAdd/
func DateAdd(adder *DateAdder) bson.Entry {
	return bson.E(operator.DateAdd, adder)
}

// DateDiff returns the difference between two dates.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/dateDiff/
func DateDiff(differ *DateDiffer) bson.Entry {
	return bson.E(operator.DateDiff, differ)
}

// DateSubtract Decrements a `Date()` object by a specified number of time units.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/dateSubtract/
func DateSubtract(subtracter *DateSubtracter) bson.Entry {
	return bson.E(operator.DateSubtract, subtracter)
}

// Divide divides one number by another and returns the result. Pass the arguments to `$divide` in an array.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/divide/
func Divide(expr1, expr2 any) bson.Entry {
	return bson.E(operator.Divide, bson.A(expr1, expr2))
}

// DocumentNumber returns the position of a document (known as the document number) in the $setWindowFields stage
// partition.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/documentNumber/
func DocumentNumber() bson.Entry {
	return bson.E(operator.DocumentNumber, bson.M())
}

// Eq Compares two values and returns:
//   - true when the values are equivalent.
//   - false when the values are not equivalent.
//
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/eq/
func Eq(expr1, expr2 any) bson.Entry {
	return bson.E(operator.Eq, bson.A(expr1, expr2))
}

// Filter selects a subset of an array to return based on the specified condition. Returns an array with only those
// elements that match the condition. The returned elements are in the original order.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/filter/
func Filter(source *FilterSource) bson.Entry {
	return bson.E(operator.Filter, source)
}

func Sum(expression any) bson.Entry {
	return bson.E(operator.Sum, expression)
}

func Subtract(expr1, expr2 any) bson.Entry {
	return bson.E(operator.Subtract, bson.A(expr1, expr2))
}
