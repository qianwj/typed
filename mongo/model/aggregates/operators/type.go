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

// Convert converts a value to a specified type.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/convert/
func Convert(c *Converter) bson.Entry {
	return bson.E(operator.Convert, c)
}

// IsNumber checks if the specified expression resolves to one of the following numeric BSON types:
//   - Integer
//   - Decimal
//   - Double
//   - Long
//
// IsNumber returns:
//   - true if the expression resolves to a number.
//   - false if the expression resolves to any other BSON type, null, or a missing field.
//
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/isNumber/
func IsNumber(expr any) bson.Entry {
	return bson.E(operator.IsNumber, expr)
}

// ToDecimal converts a value to a decimal. If the value cannot be converted to a decimal, `$toDecimal` errors. If the
// value is null or missing, `$toDecimal` returns null.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/toDecimal/
func ToDecimal(expr any) bson.Entry {
	return bson.E(operator.ToDecimal, expr)
}

// ToDouble converts a value to a double. If the value cannot be converted to an double, `$toDouble` errors. If the
// value is null or missing, `$toDouble` returns null.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/toDouble/
func ToDouble(expr any) bson.Entry {
	return bson.E(operator.ToDouble, expr)
}

// ToObjectId converts a value to an ObjectId(). If the value cannot be converted to an ObjectId, `$toObjectId` errors.
// If the value is null or missing, `$toObjectId` returns null.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/lookup/
func ToObjectId(expr any) bson.Entry {
	return bson.E(operator.ToObjectID, expr)
}