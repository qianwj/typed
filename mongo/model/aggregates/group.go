package aggregates

//type GroupId interface {
//	isGroupId()
//}
//
//type StringGroupId string
//
//func (s StringGroupId) isGroupId() {}
//
//type ComplexGroupId bson.D
//
//func (c ComplexGroupId) isGroupId() {}
//
//type GroupField struct {
//	name        string
//	accumulator string
//	expression  Expression
//}
//
//func NewGroupField(name, accumulator string, expression Expression) *GroupField {
//	return &GroupField{
//		name:        name,
//		accumulator: accumulator,
//		expression:  expression,
//	}
//}
//
//func (g *GroupField) Tag() {}
//
//func (g *GroupField) Marshal() primitive.E {
//	return primitive.E{
//		Key: g.name,
//		Value: bson.E{
//			Key:   g.accumulator,
//			Value: g.expression,
//		},
//	}
//}
