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

// DocumentNumber returns the position of a document (known as the document number) in the $setWindowFields stage
// partition.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/documentNumber/
func DocumentNumber() bson.Entry {
	return bson.E(operator.DocumentNumber, bson.M())
}

func accMulti(field, op string, exprs ...any) bson.Entry {
	if len(exprs) == 1 {
		return bson.E(field, bson.E(op, exprs[0]))
	}
	return bson.E(field, bson.E(op, bson.A(exprs...)))
}

func computeBoth(op string, expr1, expr2 any) bson.Entry {
	return bson.E(op, bson.A(expr1, expr2))
}

func computeDateWithZone(op string, date any, timezone ...any) bson.Entry {
	if len(timezone) == 0 {
		return bson.E(op, date)
	}
	return bson.E(op, bson.M(
		bson.E("date", date),
		bson.E("timezone", timezone[0]),
	))
}
