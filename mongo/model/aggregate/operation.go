package aggregate

import (
	"github.com/qianwj/typed/mongo/model"
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
