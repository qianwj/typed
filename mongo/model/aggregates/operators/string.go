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

// Concat concatenates strings and returns the concatenated string.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/concat/
func Concat(exprs ...any) bson.Entry {
	return bson.E(operator.Concat, bson.A(exprs...))
}

// Split divides a string into an array of substrings based on a delimiter. `$split` removes the delimiter and returns
// the resulting substrings as elements of an array. If the delimiter is not found in the string, $split returns the
// original string as the only element of an array.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/split/
func Split(expr any, delimiter string) bson.Entry {
	return bson.E(operator.Split, bson.A(expr, delimiter))
}
