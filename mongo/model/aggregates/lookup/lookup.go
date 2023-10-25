package lookup

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JoinCondition struct {
	From         string        `bson:"from,omitempty"`         // <collection to join>
	LocalField   string        `bson:"localField,omitempty"`   // <field from the input documents>,
	ForeignField string        `bson:"foreignField,omitempty"` // <field from the documents of the "from" collection>,
	Let          primitive.M   `bson:"let,omitempty"`          // { <var_1>: <expression>, …, <var_n>: <expression> }
	Pipeline     []primitive.D `bson:"pipeline,omitempty"`     // [ <pipeline to run on joined collection> ]
	As           string        `bson:"as,omitempty"`           // <output array field>
}

func New(from, as string) *JoinCondition {
	return &JoinCondition{
		From: from,
		As:   as,
	}
}

func (j *JoinCondition) Join(localField, foreignField string) *JoinCondition {
	j.LocalField = localField
	j.ForeignField = foreignField
	return j
}

func (j *JoinCondition) Pipe(pipeline []primitive.D, let primitive.M) *JoinCondition {
	j.Let = let
	j.Pipeline = pipeline
	return j
}

//
//func (l *Lookup) Marshal() primitive.D {
//	res := primitive.D{
//		{Key: "from", Value: l.from},
//		{Key: "as", Value: l.as},
//	}
//	if l.localField != "" && l.foreignField != "" {
//		res = append(
//			res,
//			primitive.E{Key: "localField", Value: l.localField},
//			primitive.E{Key: "foreignField", Value: l.foreignField},
//		)
//	}
//	if len(l.let) > 0 {
//		res = append(res, primitive.E{Key: "let", Value: l.let})
//	}
//	if l.pipeline != nil {
//		res = append(res, primitive.E{Key: "pipeline", Value: l.pipeline.Marshal()})
//	}
//	return res
//}
//
//func (l *Lookup) Tag() {}
//
//func (l *Lookup) ToMap() primitive.M {
//	arr, res := l.Marshal(), primitive.M{}
//	for _, e := range arr {
//		res[e.Key] = e.Value
//	}
//	return res
//}
//
//type GraphLookup struct {
//	from                    string          // <collection>,
//	startWith               string          // <expression>,
//	connectFromField        string          // <string>,
//	connectToField          string          // <string>,
//	as                      string          // <string>,
//	maxDepth                int             // <number>,
//	depthField              string          // <string>,
//	restrictSearchWithMatch *filters.Filter // <document>
//}
//
//func (l *GraphLookup) Tag() {}
//
//func (l *GraphLookup) Marshal() primitive.D {
//	// complete this function
//	return primitive.D{}
//}
//
//func (l *GraphLookup) ToMap() primitive.M {
//	return primitive.M{}
//}
