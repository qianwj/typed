package aggregates

type Variable interface {
	varTag()
}

type StringVariable string

func (v StringVariable) varTag() {}

func (v StringVariable) IsExpression() {}

type LongVariable int64

func (v LongVariable) varTag() {}
