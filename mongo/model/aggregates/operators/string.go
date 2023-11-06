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
	"github.com/qianwj/typed/mongo/model/regex"
	"github.com/qianwj/typed/mongo/operator"
	"github.com/qianwj/typed/mongo/util"
	"strings"
)

// Concat concatenates strings and returns the concatenated string.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/concat/
func Concat(exprs ...any) bson.Entry {
	return bson.E(operator.Concat, bson.A(exprs...))
}

// IndexOfBytes searches a string for an occurrence of a substring and returns the UTF-8 byte index (zero-based) of the
// first occurrence. If the substring is not found, returns -1.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/indexOfBytes/
func IndexOfBytes(str, substr any, ranges ...any) bson.Entry {
	return bson.E(operator.IndexOfBytes, bson.A(str, substr).Append(ranges...))
}

// IndexOfCP searches a string for an occurrence of a substring and returns the UTF-8 code point index (zero-based) of
// the first occurrence. If the substring is not found, returns -1.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/indexOfCP/
func IndexOfCP(str, substr any, ranges ...any) bson.Entry {
	return bson.E(operator.IndexOfCP, bson.A(str, substr).Append(ranges...))
}

// Ltrim removes whitespace characters, including null, or the specified characters from the beginning of a string.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/ltrim/
func Ltrim(input, chars any) bson.Entry {
	return bson.E(operator.Ltrim, bson.M(
		bson.E("inputs", input),
		bson.E("chars", chars),
	))
}

// RegexFind provides regular expression (regex) pattern matching capability in aggregation expressions. If a match is
// found, returns a document that contains information on the first match. If a match is not found, returns null.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/regexFind/
func RegexFind(input, regexExp any, opts ...regex.Options) bson.Entry {
	val := bson.M(
		bson.E("input", input),
		bson.E("regex", regexExp),
	)
	if len(opts) > 0 {
		val["options"] = strings.Join(util.Map(opts, regex.Options.String), "")
	}
	return bson.E(operator.RegexFind, val)
}

// Rtrim removes whitespace characters, including null, or the specified characters from the end of a string.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/rtrim/
func Rtrim(input, chars any) bson.Entry {
	return bson.E(operator.Rtrim, bson.M(
		bson.E("inputs", input),
		bson.E("chars", chars),
	))
}

// Split divides a string into an array of substrings based on a delimiter. `$split` removes the delimiter and returns
// the resulting substrings as elements of an array. If the delimiter is not found in the string, $split returns the
// original string as the only element of an array.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/split/
func Split(expr any, delimiter string) bson.Entry {
	return bson.E(operator.Split, bson.A(expr, delimiter))
}

// ToLower converts a string to lowercase, returning the result.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/toLower/
func ToLower(expr any) bson.Entry {
	return bson.E(operator.ToLower, expr)
}

// ToUpper converts a string to uppercase, returning the result.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/toUpper/
func ToUpper(expr any) bson.Entry {
	return bson.E(operator.ToUpper, expr)
}

// Trim removes whitespace characters, including null, or the specified characters from the beginning and end of a
// string.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/trim/
func Trim(input, chars any) bson.Entry {
	return bson.E(operator.Trim, bson.M(
		bson.E("inputs", input),
		bson.E("chars", chars),
	))
}
