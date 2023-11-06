package not

import (
	"github.com/qianwj/typed/mongo/bson"
	"github.com/qianwj/typed/mongo/model/regex"
	"github.com/qianwj/typed/mongo/operator"
	rawbson "go.mongodb.org/mongo-driver/bson"
)

type OperatorExpression struct {
	data   bson.UnorderedMap
	regexp *regex.Matcher
	val    any
}

func NewOperatorExpression() *OperatorExpression {
	return &OperatorExpression{data: bson.M()}
}

func Like(regexp *regex.Matcher) *OperatorExpression {
	return NewOperatorExpression().Like(regexp)
}

func Gt(val any) *OperatorExpression {
	return NewOperatorExpression().Gt(val)
}

func Gte(val any) *OperatorExpression {
	return NewOperatorExpression().Gte(val)
}

func Lt(val any) *OperatorExpression {
	return NewOperatorExpression().Lt(val)
}

func Lte(val any) *OperatorExpression {
	return NewOperatorExpression().Lte(val)
}

func Eq(val any) *OperatorExpression {
	return NewOperatorExpression().Eq(val)
}

func (n *OperatorExpression) Like(regexp *regex.Matcher) *OperatorExpression {
	n.regexp = regexp
	clear(n.data)
	return n
}

func (n *OperatorExpression) Gt(val any) *OperatorExpression {
	n.data[operator.Gt] = val
	return n
}

func (n *OperatorExpression) Gte(val any) *OperatorExpression {
	n.data[operator.Gte] = val
	return n
}

func (n *OperatorExpression) Lt(val any) *OperatorExpression {
	n.data[operator.Lt] = val
	return n
}

func (n *OperatorExpression) Lte(val any) *OperatorExpression {
	n.data[operator.Lte] = val
	return n
}

func (n *OperatorExpression) Eq(val any) *OperatorExpression {
	n.val = val
	clear(n.data)
	return n
}

func (n *OperatorExpression) MarshalBSON() ([]byte, error) {
	if len(n.data) > 0 {
		return rawbson.Marshal(n.data)
	}
	if n.regexp != nil {
		return n.regexp.MarshalBSON()
	}
	return rawbson.Marshal(n.val)
}
