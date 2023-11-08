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

type AccumulatorOptions struct {
	Init           string `bson:"init"`               // <code>,
	InitArgs       []any  `bson:"initArgs,omitempty"` // <array expression>,        // Optional
	Accumulate     string `bson:"accumulate"`         // <code>,
	AccumulateArgs []any  `bson:"accumulateArgs"`     // <array expression>,
	Merge          string `bson:"merge"`              // <code>,
	Finalize       string `bson:"finalize,omitempty"` // <code>,                    // Optional
	Lang           string `bson:"lang"`               // <string>
}

func NewAccumulatorOptions(init, acc, merge, lang string) *AccumulatorOptions {
	return &AccumulatorOptions{
		Init:       init,
		Accumulate: acc,
		Merge:      merge,
		Lang:       lang,
	}
}

func (a *AccumulatorOptions) AddInitArgs(args ...any) *AccumulatorOptions {
	a.InitArgs = append(a.InitArgs, args...)
	return a
}

func (a *AccumulatorOptions) AddAccumulateArgs(args ...any) *AccumulatorOptions {
	a.AccumulateArgs = append(a.AccumulateArgs, args...)
	return a
}

func (a *AccumulatorOptions) SetFinalize(finalize string) *AccumulatorOptions {
	a.Finalize = finalize
	return a
}
