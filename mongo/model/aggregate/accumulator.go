package aggregate

//
//import (
//	"github.com/qianwj/typed/mongo/bson"
//	"github.com/qianwj/typed/mongo/model/filters/text"
//	"go.mongodb.org/mongo-driver/bson/primitive"
//)
//
//type Accumulate struct {
//	init           string         // <code>,
//	initArgs       bson.Array     // <array expression>,        // Optional
//	accumulate     string         // <code>,
//	accumulateArgs bson.Array     // <array expression>,
//	merge          string         // <code>,
//	finalize       *string        // <code>,                    // Optional
//	lang           text.Langeuage // <string>
//}
//
//func NewAccumulate(init, accumulate, merge string, accumulateArgs bson.Array, lang text.Langeuage) *Accumulate {
//	return &Accumulate{
//		init:           init,
//		accumulate:     accumulate,
//		accumulateArgs: accumulateArgs,
//		merge:          merge,
//		lang:           lang,
//	}
//}
//
//func (a *Accumulate) Tag() {}
//
//func (a *Accumulate) Marshal() primitive.D {
//	res := primitive.D{
//		{Key: "init", Value: a.init},
//	}
//	if a.initArgs.Size() > 0 {
//		res = append(res, primitive.E{
//			Key: "initArgs", Value: a.initArgs,
//		})
//	}
//	res = append(res,
//		primitive.E{Key: "accumulate", Value: a.accumulate},
//		primitive.E{Key: "accumulateArgs", Value: a.accumulateArgs},
//		primitive.E{Key: "merge", Value: a.merge},
//	)
//	if a.finalize != nil {
//		res = append(res, primitive.E{
//			Key: "finalize", Value: *a.finalize,
//		})
//	}
//	res = append(res,
//		primitive.E{Key: "lang", Value: a.lang},
//	)
//	return res
//}
//
//func (a *Accumulate) ToMap() primitive.M {
//	res := primitive.M{
//		"init":           a.init,
//		"accumulate":     a.accumulate,
//		"accumulateArgs": a.accumulateArgs,
//		"merge":          a.merge,
//		"lang":           a.lang,
//	}
//	if a.initArgs.Size() > 0 {
//		res["initArgs"] = a.initArgs
//	}
//	if a.finalize != nil {
//		res["finalize"] = *a.finalize
//	}
//	return res
//}
