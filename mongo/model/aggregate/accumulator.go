package aggregate

import (
	"github.com/qianwj/typed/mongo/model/filter/text"
	"go.mongodb.org/mongo-driver/bson"
)

type Accumulate struct {
	init           string         // <code>,
	initArgs       bson.A         // <array expression>,        // Optional
	accumulate     string         // <code>,
	accumulateArgs bson.A         // <array expression>,
	merge          string         // <code>,
	finalize       *string        // <code>,                    // Optional
	lang           text.Langeuage // <string>
}

func NewAccumulate(init, accumulate, merge string, accumulateArgs bson.A, lang text.Langeuage) *Accumulate {
	return &Accumulate{
		init:           init,
		accumulate:     accumulate,
		accumulateArgs: accumulateArgs,
		merge:          merge,
		lang:           lang,
	}
}

func (a *Accumulate) Marshal() bson.D {
	res := bson.D{
		{Key: "init", Value: a.init},
	}
	if len(a.initArgs) > 0 {
		res = append(res, bson.E{
			Key: "initArgs", Value: a.initArgs,
		})
	}
	res = append(res,
		bson.E{Key: "accumulate", Value: a.accumulate},
		bson.E{Key: "accumulateArgs", Value: a.accumulateArgs},
		bson.E{Key: "merge", Value: a.merge},
	)
	if a.finalize != nil {
		res = append(res, bson.E{
			Key: "finalize", Value: *a.finalize,
		})
	}
	res = append(res,
		bson.E{Key: "lang", Value: a.lang},
	)
	return res
}
