package aggregate

import (
	"github.com/qianwj/typed/mongo/model"
	"github.com/qianwj/typed/mongo/model/filter"
	"github.com/qianwj/typed/mongo/model/operator"
	"go.mongodb.org/mongo-driver/bson"
)

func Abs[E model.Bson](expression E) bson.D {
	return bson.D{
		{Key: operator.Abs, Value: expression},
	}
}

func Accumulator(acc *Accumulate) bson.D {
	return bson.D{
		{Key: operator.Accumulator, Value: acc.Marshal()},
	}
}

func Add(vars ...model.Variable) bson.D {
	return bson.D{
		{Key: operator.Add, Value: vars},
	}
}

func AddToSet(variable model.StringVariable) bson.D {
	return bson.D{
		{Key: operator.AddToSet, Value: variable},
	}
}

func And(filters ...*filter.Filter) bson.D {
	res := make(bson.A, len(filters))
	for i, f := range filters {
		res[i] = f.Marshal()
	}
	return bson.D{
		{Key: operator.And, Value: res},
	}
}

func Avg[E model.Bson](expr E) bson.D {
	return bson.D{
		{Key: operator.Avg, Value: expr},
	}
}

func Ceil(variable model.StringVariable) bson.D {
	return bson.D{
		{Key: operator.Ceil, Value: variable},
	}
}

func Compare[E model.Bson](expr1, expr2 E) bson.D {
	return bson.D{
		{Key: operator.Cmp, Value: bson.A{
			expr1, expr2,
		}},
	}
}

func Concat[E model.Bson](expr1, expr2 E) bson.D {
	return bson.D{
		{Key: operator.Concat, Value: bson.A{
			expr1, expr2,
		}},
	}
}

func Count() bson.D {
	return bson.D{{Key: operator.Count, Value: bson.M{}}}
}
