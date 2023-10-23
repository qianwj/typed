package aggregates

//import (
//	"github.com/qianwj/typed/mongo/bson"
//	"github.com/qianwj/typed/mongo/model/filters"
//	"go.mongodb.org/mongo-driver/bson/primitive"
//)
//
//type Lookup struct {
//	from         string       // <collection to join>
//	localField   string       // <field from the input documents>,
//	foreignField string       // <field from the documents of the "from" collection>,
//	let          []Expression // { <var_1>: <expression>, …, <var_n>: <expression> }
//	pipeline     bson.Array   // [ <pipeline to run on joined collection> ]
//	as           string       // <output array field>
//}
//
//func NewLookup(from, as string) *Lookup {
//	return &Lookup{
//		from: from,
//		as:   as,
//	}
//}
//
//func (l *Lookup) Join(localField, foreignField string) *Lookup {
//	l.localField = localField
//	l.foreignField = foreignField
//	l.let = nil
//	l.pipeline = nil
//	return l
//}
//
//func (l *Lookup) Pipeline(pipeline bson.Array, let ...Expression) *Lookup {
//	l.localField = ""
//	l.foreignField = ""
//	l.let = let
//	l.pipeline = pipeline
//	return l
//}
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
