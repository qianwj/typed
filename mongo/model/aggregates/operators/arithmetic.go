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
	"github.com/qianwj/typed/mongo/operator"
)

// Abs returns the absolute value of a number.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/abs/
func Abs[E expressions.Number](expr E) bson.Entry {
	return bson.E(operator.Abs, expr)
}

// Add adds numbers together or adds numbers and a date. If one of the arguments is a date, `$add` treats the other
// arguments as milliseconds to add to the date.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/add/
func Add(exprs ...any) bson.Entry {
	return bson.E(operator.Add, bson.A(exprs...))
}

// Ceil returns the smallest integer greater than or equal to the specified number.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/ceil/
func Ceil[E expressions.Number](expr E) bson.Entry {
	return bson.E(operator.Ceil, expr)
}

// Divide divides one number by another and returns the result. Pass the arguments to `$divide` in an array.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/divide/
func Divide[T, R expressions.Number](expr1 T, expr2 R) bson.Entry {
	return computeBoth(operator.Divide, expr1, expr2)
}

// Exp raises Euler's number (i.e. e ) to the specified exponent and returns the result.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/exp/
func Exp[E expressions.Number](exponent E) bson.Entry {
	return bson.E(operator.Exp, exponent)
}

// Floor returns the largest integer less than or equal to the specified number.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/floor/
func Floor[E expressions.Number](expr E) bson.Entry {
	return bson.E(operator.Floor, expr)
}

// Ln calculates the natural logarithm ln (i.e log e) of a number and returns the result as a double.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/ln/
func Ln[E expressions.Number](expr E) bson.Entry {
	return bson.E(operator.Ln, expr)
}

// Log calculates the log of a number in the specified base and returns the result as a double.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/log/
func Log[T, R expressions.Number](number T, base R) bson.Entry {
	return bson.E(operator.Log, bson.A(number, base))
}

// Lg calculates the log base 10 of a number and returns the result as a double.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/log10/
func Lg[E expressions.Number](expr E) bson.Entry {
	return bson.E(operator.Log10, expr)
}

// Mod divides one number by another and returns the remainder.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/mod/
func Mod[T, R expressions.Number](expr1 T, expr2 R) bson.Entry {
	return computeBoth(operator.Mod, expr1, expr2)
}

// Multiply Multiplies numbers together and returns the result. Pass the arguments to `$multiply` in an array.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/multiply/
func Multiply(exprs ...any) bson.Entry {
	return bson.E(operator.Multiply, bson.A(exprs...))
}

// Pow raises a number to the specified exponent and returns the result.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/pow/
func Pow[T, R expressions.Number](number T, exponent R) bson.Entry {
	return bson.E(operator.Pow, bson.A(number, exponent))
}

// Round rounds a number to a whole integer or to a specified decimal place.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/round/
func Round(number any, place ...any) bson.Entry {
	return bson.E(operator.Round, bson.A(number).Append(place...))
}

// Sqrt calculates the square root of a positive number and returns the result as a double.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/sqrt/
func Sqrt[E expressions.Number](expr E) bson.Entry {
	return bson.E(operator.Sqrt, expr)
}

// Subtract Subtracts two numbers to return the difference, or two dates to return the difference in milliseconds, or a
// date and a number in milliseconds to return the resulting date.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/subtract/
func Subtract[T, R expressions.Number](expr1 T, expr2 R) bson.Entry {
	return computeBoth(operator.Subtract, expr1, expr2)
}

// Trunc truncates a number to a whole integer or to a specified decimal place.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/trunc/
func Trunc(number any, place ...any) bson.Entry {
	return bson.E(operator.Trunc, bson.A(number).Append(place...))
}
