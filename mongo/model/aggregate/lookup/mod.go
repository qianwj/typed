package lookup

import (
	"github.com/qianwj/typed/mongo/model/aggregate/pipe"
	"go.mongodb.org/mongo-driver/bson"
)

type Lookup struct {
	from         string         // <collection to join>
	localField   string         // <field from the input documents>,
	foreignField string         // <field from the documents of the "from" collection>,
	let          bson.M         // { <var_1>: <expression>, …, <var_n>: <expression> }
	pipeline     *pipe.Pipeline // [ <pipeline to run on joined collection> ]
	as           string         // <output array field>
}

func New(from, as string) *Lookup {
	return &Lookup{
		from: from,
		as:   as,
	}
}

func (l *Lookup) Join(localField, foreignField string) *Lookup {
	l.localField = localField
	l.foreignField = foreignField
	l.let = nil
	l.pipeline = nil
	return l
}

func (l *Lookup) MultiJoin(let bson.M, pipeline *pipe.Pipeline) *Lookup {
	l.localField = ""
	l.foreignField = ""
	l.let = let
	l.pipeline = pipeline
	return l
}

func (l *Lookup) Marshal() bson.D {
	res := bson.D{
		{Key: "from", Value: l.from},
		{Key: "as", Value: l.as},
	}
	if l.localField != "" && l.foreignField != "" {
		res = append(
			res,
			bson.E{Key: "localField", Value: l.localField},
			bson.E{Key: "foreignField", Value: l.foreignField},
		)
	}
	if len(l.let) > 0 {
		res = append(res, bson.E{Key: "let", Value: l.let})
	}
	if l.pipeline != nil {
		res = append(res, bson.E{Key: "pipeline", Value: l.pipeline.Marshal()})
	}
	return res
}
