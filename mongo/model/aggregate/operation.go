package aggregate

//
//import (
//	"github.com/qianwj/typed/mongo/model"
//	"github.com/qianwj/typed/mongo/model/filters"
//	operator2 "github.com/qianwj/typed/mongo/operator"
//	"go.mongodb.org/mongo-driver/bson"
//)
//
//func Abs[E model.Bson](expression E) bson.D {
//	return bson.D{
//		{Key: operator2.Abs, Value: expression},
//	}
//}
//
//func Accumulator(acc *Accumulate) bson.D {
//	return bson.D{
//		{Key: operator2.Accumulator, Value: acc.Marshal()},
//	}
//}
//
//func Add(vars ...model.Variable) bson.D {
//	return bson.D{
//		{Key: operator2.Add, Value: vars},
//	}
//}
//
//func AddToSet(variable model.StringVariable) bson.D {
//	return bson.D{
//		{Key: operator2.AddToSet, Value: variable},
//	}
//}
//
//func And(filters ...*filters.Filter) bson.D {
//	res := make(bson.A, len(filters))
//	for i, f := range filters {
//		res[i] = f.Marshal()
//	}
//	return bson.D{
//		{Key: operator2.And, Value: res},
//	}
//}
//
//func Avg[E model.Bson](expr E) bson.D {
//	return bson.D{
//		{Key: operator2.Avg, Value: expr},
//	}
//}
//
//func Ceil(variable model.StringVariable) bson.D {
//	return bson.D{
//		{Key: operator2.Ceil, Value: variable},
//	}
//}
//
//func Compare[E model.Bson](expr1, expr2 E) bson.D {
//	return bson.D{
//		{Key: operator2.Cmp, Value: bson.A{
//			expr1, expr2,
//		}},
//	}
//}
//
//func Concat[E model.Bson](expr1, expr2 E) bson.D {
//	return bson.D{
//		{Key: operator2.Concat, Value: bson.A{
//			expr1, expr2,
//		}},
//	}
//}
//
//func Count() bson.D {
//	return bson.D{{Key: operator2.Count, Value: bson.M{}}}
//}
