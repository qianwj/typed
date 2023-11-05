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

// Cmp Compares two values and returns:
//   - -1 if the first value is less than the second.
//   - 1 if the first value is greater than the second.
//   - 0 if the two values are equivalent.
//
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/cmp/
func Cmp(expr1, expr2 any) bson.Entry {
	return computeBoth(operator.Cmp, expr1, expr2)
}

// Divide divides one number by another and returns the result. Pass the arguments to `$divide` in an array.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/divide/
func Divide(expr1, expr2 any) bson.Entry {
	return computeBoth(operator.Divide, expr1, expr2)
}

// Eq Compares two values and returns:
//   - true when the values are equivalent.
//   - false when the values are not equivalent.
//
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/eq/
func Eq(expr1, expr2 any) bson.Entry {
	return computeBoth(operator.Eq, expr1, expr2)
}

// Gt Compares two values and returns:
//   - true when the first value is greater than the second value.
//   - false when the first value is less than or equivalent to the second value.
//
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/gt/
func Gt(expr1, expr2 any) bson.Entry {
	return computeBoth(operator.Gt, expr1, expr2)
}

// Gte Compares two values and returns:
//   - true when the first value is greater than or equivalent the second value.
//   - false when the first value is less than or equivalent to the second value.
//
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/gte/
func Gte(expr1, expr2 any) bson.Entry {
	return computeBoth(operator.Gte, expr1, expr2)
}

// In returns a boolean indicating whether a specified value is in an array.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/in/
func In(elemExpr, arrExpr any) bson.Entry {
	return computeBoth(operator.In, elemExpr, arrExpr)
}

// Lt Compares two values and returns:
//   - true when the first value is less than the second value.
//   - false when the first value is greater than or equivalent to the second value.
//
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/lt/
func Lt(expr1, expr2 any) bson.Entry {
	return computeBoth(operator.Lt, expr1, expr2)
}

// Lte Compares two values and returns:
//   - true when the first value is less than or equivalent to the second value.
//   - false when the first value is greater than the second value.
//
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/lte/
func Lte(expr1, expr2 any) bson.Entry {
	return computeBoth(operator.Lte, expr1, expr2)
}

// Subtract Subtracts two numbers to return the difference, or two dates to return the difference in milliseconds, or a
// date and a number in milliseconds to return the resulting date.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/subtract/
func Subtract(expr1, expr2 any) bson.Entry {
	return computeBoth(operator.Subtract, expr1, expr2)
}

func computeBoth(op string, expr1, expr2 any) bson.Entry {
	return bson.E(op, bson.A(expr1, expr2))
}
