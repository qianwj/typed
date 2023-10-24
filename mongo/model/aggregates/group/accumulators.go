package group

import (
	"github.com/qianwj/typed/mongo/operator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Accumulator primitive.E

func AddToSet(field string, expression any) Accumulator {
	return newAccumulator(field, operator.AddToSet, expression)
}

func Avg(field string, expression any) Accumulator {
	return newAccumulator(field, operator.Avg, expression)
}

func Count(field string, expression any) Accumulator {
	return newAccumulator(field, operator.Count, expression)
}

func First(field string, expression any) Accumulator {
	return newAccumulator(field, operator.First, expression)
}

func Max(field string, expression any) Accumulator {
	return newAccumulator(field, operator.Max, expression)
}

func Push(field string, expression any) Accumulator {
	return newAccumulator(field, operator.Push, expression)
}

func Sum(field string, expression any) Accumulator {
	return newAccumulator(field, operator.Sum, expression)
}

func newAccumulator(field, operator string, expression any) Accumulator {
	return Accumulator{
		Key: field,
		Value: primitive.D{
			{Key: operator, Value: expression},
		},
	}
}
