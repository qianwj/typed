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
	"github.com/qianwj/typed/mongo/model/sorts"
	"github.com/qianwj/typed/mongo/operator"
	"time"
)

// Abs returns the absolute value of a number.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/abs/
func Abs(numOrExpr any) bson.Entry {
	return bson.E(operator.Abs, numOrExpr)
}

// Accumulator defines a custom accumulator operator. Accumulators are operators that maintain their state (e.g. totals,
// maximums, minimums, and related data) as documents progress through the pipeline. Use the `$accumulator` operator to
// execute your own JavaScript functions to implement behavior not supported by the MongoDB Query Language.
// See also $function.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/accumulator/
func Accumulator(acc *AccumulatorSource) bson.Entry {
	return bson.E(operator.Accumulator, acc)
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

// Bottom returns the bottom element within a group according to the specified sort order.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/bottom/
func Bottom(sortBy *sorts.Options, output any) bson.Entry {
	return bson.E(operator.Bottom, bson.M(
		bson.E("sortBy", sortBy),
		bson.E("output", output),
	))
}

// Ceil returns the smallest integer greater than or equal to the specified number.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/ceil/
func Ceil(numOrExpr any) bson.Entry {
	return bson.E(operator.Ceil, numOrExpr)
}

// Concat concatenates strings and returns the concatenated string.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/concat/
func Concat(exprs ...any) bson.Entry {
	return bson.E(operator.Concat, bson.A(exprs...))
}

// ConcatArrays concatenates arrays to return the concatenated array.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/concatArrays/
func ConcatArrays(exprs ...any) bson.Entry {
	return bson.E(operator.ConcatArrays, bson.A(exprs...))
}

// Cond evaluates a boolean expression to return one of the two specified return expressions.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/cond/
func Cond(assuming, thenCase, elseCase any) bson.Entry {
	return bson.E(operator.Cond, bson.A(assuming, thenCase, elseCase))
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

// DocumentNumber returns the position of a document (known as the document number) in the $setWindowFields stage
// partition.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/documentNumber/
func DocumentNumber() bson.Entry {
	return bson.E(operator.DocumentNumber, bson.M())
}

// Filter selects a subset of an array to return based on the specified condition. Returns an array with only those
// elements that match the condition. The returned elements are in the original order.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/filter/
func Filter(source *FilterSource) bson.Entry {
	return bson.E(operator.Filter, source)
}

// First returns the result of an expression for the first document in a group of documents. Only meaningful when
// documents are in a defined order.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/first/
func First(expr any) bson.Entry {
	return bson.E(operator.First, expr)
}

// Floor returns the largest integer less than or equal to the specified number.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/floor/
func Floor(expr any) bson.Entry {
	return bson.E(operator.Floor, expr)
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

// Hour returns the hour portion of a date as a number between 0 and 23.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/hour/
func Hour(date time.Time, timezone ...string) bson.Entry {
	if len(timezone) == 0 {
		return bson.E(operator.Hour, date)
	}
	return bson.E(operator.Hour, bson.M(
		bson.E("date", date),
		bson.E("timezone", timezone[0]),
	))
}

// IfNull evaluates input expressions for null values and returns:
//   - The first non-null input expression value found.
//   - A replacement expression value if all input expressions evaluate to null.
//
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/ifNull/
func IfNull(replacement any, inputs ...any) bson.Entry {
	return bson.E(operator.IfNull, bson.A(inputs...).Append(replacement))
}

// Last returns the result of an expression for the last document in a group of documents. Only meaningful when
// documents are in a defined order.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/last/
func Last(expr any) bson.Entry {
	return bson.E(operator.Last, expr)
}

func Sum(expr any) bson.Entry {
	return bson.E(operator.Sum, expr)
}
