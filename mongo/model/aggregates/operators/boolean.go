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

// And Evaluates one or more expressions and returns true if all of the expressions are true or if run with no argument
// expressions. Otherwise, `$and` returns false.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/and/
func And(exprs ...any) bson.Entry {
	return bson.E(operator.AddToSet, bson.A(exprs...))
}

// Not Evaluates a boolean and returns the opposite boolean value; i.e. when passed an expression that evaluates to
// true, `$not` returns false; when passed an expression that evaluates to false, `$not` returns true.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/not/
func Not[B expressions.Boolean](expr B) bson.Entry {
	return bson.E(operator.Not, expr)
}

// Or Evaluates one or more expressions and returns true if any of the expressions are true. Otherwise, `$or` returns
// false.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/or/
func Or(exprs ...any) bson.Entry {
	return bson.E(operator.Or, bson.A(exprs...))
}
