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
	"github.com/qianwj/typed/mongo/util"
)

// Cond evaluates a boolean expression to return one of the two specified return expressions.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/cond/
func Cond(assuming, thenCase, elseCase any) bson.Entry {
	return bson.E(operator.Cond, bson.A(assuming, thenCase, elseCase))
}

// IfNull evaluates input expressions for null values and returns:
//   - The first non-null input expression value found.
//   - A replacement expression value if all input expressions evaluate to null.
//
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/ifNull/
func IfNull(replacement any, inputs ...any) bson.Entry {
	return bson.E(operator.IfNull, bson.A(inputs...).Append(replacement))
}

// Switch evaluates a series of case expressions. When it finds an expression which evaluates to true, `$switch`
// executes a specified expression and breaks out of the control flow.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/switch/
func Switch(cond *SwitchCondition) bson.Entry {
	return bson.E(operator.Switch, cond)
}

type SwitchCondition struct {
	branches    []*Branch
	defaultCase any
}

func NewSwitch(defaultCase any) *SwitchCondition {
	return &SwitchCondition{
		defaultCase: defaultCase,
	}
}

func (s *SwitchCondition) AddBranch(switchCase, then any) *SwitchCondition {
	s.branches = append(s.branches, &Branch{
		switchCase: switchCase,
		then:       then,
	})
	return s
}

func (s *SwitchCondition) MarshalBSON() ([]byte, error) {
	return bson.M(
		bson.E("branches", bson.A(util.ToAny(s.branches)...)),
		bson.E("default", s.defaultCase),
	).Marshal()
}

type Branch struct {
	switchCase any
	then       any
}

func (b *Branch) MarshalBSON() ([]byte, error) {
	return bson.M(
		bson.E("case", b.switchCase),
		bson.E("then", b.then),
	).Marshal()
}
