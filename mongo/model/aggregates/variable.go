package aggregates

type Variable interface {
	IsVariable()
}

type StringVariable string

func (v StringVariable) IsVariable() {}

func (v StringVariable) IsExpression() {}

type LongVariable int64

func (v LongVariable) varTag() {}
