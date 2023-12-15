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

package options

import "go.mongodb.org/mongo-driver/mongo/options"

type Collation struct {
	internal options.Collation
}

func NewCollation() *Collation {
	return &Collation{}
}

func (c *Collation) Locale(locale string) *Collation {
	c.internal.Locale = locale
	return c
}

func (c *Collation) CaseLevel(caseLevel bool) *Collation {
	c.internal.CaseLevel = caseLevel
	return c
}

func (c *Collation) CaseFirst(caseFirst string) *Collation {
	c.internal.CaseFirst = caseFirst
	return c
}

func (c *Collation) Strength(strength int) *Collation {
	c.internal.Strength = strength
	return c
}

func (c *Collation) NumericOrdering(numericOrdering bool) *Collation {
	c.internal.NumericOrdering = numericOrdering
	return c
}

func (c *Collation) Alternate(alternate string) *Collation {
	c.internal.Alternate = alternate
	return c
}

func (c *Collation) MaxVariable(maxVariable string) *Collation {
	c.internal.MaxVariable = maxVariable
	return c
}

func (c *Collation) Normalization(normalization bool) *Collation {
	c.internal.Normalization = normalization
	return c
}

func (c *Collation) Backwards(backwards bool) *Collation {
	c.internal.Backwards = backwards
	return c
}

func (c *Collation) Raw() *options.Collation {
	return &(c.internal)
}
