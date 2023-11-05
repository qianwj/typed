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

// Ceil returns the smallest integer greater than or equal to the specified number.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/ceil/
func Ceil(numOrExpr any) bson.Entry {
	return bson.E(operator.Ceil, numOrExpr)
}

// Divide divides one number by another and returns the result. Pass the arguments to `$divide` in an array.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/divide/
func Divide(expr1, expr2 any) bson.Entry {
	return computeBoth(operator.Divide, expr1, expr2)
}

// Floor returns the largest integer less than or equal to the specified number.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/floor/
func Floor(expr any) bson.Entry {
	return bson.E(operator.Floor, expr)
}

// Mod divides one number by another and returns the remainder.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/mod/
func Mod(expr1, expr2 any) bson.Entry {
	return computeBoth(operator.Mod, expr1, expr2)
}

// Multiply Multiplies numbers together and returns the result. Pass the arguments to `$multiply` in an array.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/multiply/
func Multiply(exprs ...any) bson.Entry {
	return bson.E(operator.Multiply, bson.A(exprs...))
}

// Subtract Subtracts two numbers to return the difference, or two dates to return the difference in milliseconds, or a
// date and a number in milliseconds to return the resulting date.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/subtract/
func Subtract(expr1, expr2 any) bson.Entry {
	return computeBoth(operator.Subtract, expr1, expr2)
}
