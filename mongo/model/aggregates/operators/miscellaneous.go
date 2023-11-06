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

// GetField returns the value of a specified field from a document. If you don't specify an object, `$getField`
// returns the value of the field from $$CURRENT.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/getField/
func GetField(field, input any) bson.Entry {
	return bson.E(operator.GetField, bson.M(bson.E("field", field), bson.E("input", input)))
}

// Rand returns a random float between 0 and 1 each time it is called.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/rand/
func Rand() bson.Entry {
	return bson.E(operator.Rand, bson.M())
}

// SampleRate matches a random selection of input documents. The number of documents selected approximates the sample
// rate expressed as a percentage of the total number of documents.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/sampleRate/
func SampleRate[F bson.Float](rate F) bson.Entry {
	return bson.E(operator.SampleRate, rate)
}

// ToHashedIndexKey computes and returns the hash value of the input expression using the same hash function that MongoDB uses to create
// a hashed index. A hash function maps a key or string to a fixed-size numeric value.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/toHashedIndexKey/
func ToHashedIndexKey(strToHash any) bson.Entry {
	return bson.E(operator.ToHashedIndexKey, strToHash)
}
