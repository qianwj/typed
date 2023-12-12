/** MIT License
 *
 * Copyright (c) 2022-2024 qianwj
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package aggregates

import (
	"github.com/qianwj/typed/mongo/util"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

type UnwindOptions struct {
	path                       string
	includeArrayIndex          *string
	preserveNullAndEmptyArrays *bool
}

func NewUnwindOptions(path string) *UnwindOptions {
	return &UnwindOptions{path: path}
}

func (u *UnwindOptions) IncludeArrayIndex(includeArrayIndex string) *UnwindOptions {
	u.includeArrayIndex = util.ToPtr(includeArrayIndex)
	return u
}

func (u *UnwindOptions) PreserveNullAndEmptyArrays(preserveNullAndEmptyArrays bool) *UnwindOptions {
	u.preserveNullAndEmptyArrays = util.ToPtr(preserveNullAndEmptyArrays)
	return u
}

func (u *UnwindOptions) MarshalBSON() ([]byte, error) {
	idx, doc := bsoncore.AppendDocumentStart(nil)
	doc = bsoncore.AppendStringElement(doc, "path", u.path)
	if util.IsNonNil(u.includeArrayIndex) {
		doc = bsoncore.AppendStringElement(doc, "includeArrayIndex", util.FromPtr(u.includeArrayIndex))
	}
	if util.IsNonNil(u.preserveNullAndEmptyArrays) {
		doc = bsoncore.AppendBooleanElement(doc, "preserveNullAndEmptyArrays", util.FromPtr(u.preserveNullAndEmptyArrays))
	}
	return bsoncore.AppendDocumentEnd(doc, idx)
}
