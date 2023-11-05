package operators

import (
	"github.com/qianwj/typed/mongo/bson"
	"github.com/qianwj/typed/mongo/operator"
)

// FirstN returns an aggregation of the first n elements within a group. The elements returned are meaningful only if
// in a specified sort order. If the group contains fewer than n elements, `$firstN` returns all elements in the group.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/firstN/
func FirstN(input, n any) bson.Entry {
	return bson.E(operator.FirstN, bson.D(
		bson.E("input", input),
		bson.E("n", n),
	))
}

// ArrayElementFirstN returns a specified number of elements from the beginning of an array.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/firstN-array-element/
func ArrayElementFirstN(n, input any) bson.Entry {
	return bson.E(operator.FirstN, bson.D(
		bson.E("n", n),
		bson.E("input", input),
	))
}

// LastN returns an aggregation of the last n elements within a group. The elements returned are meaningful only if in
// a specified sort order. If the group contains fewer than n elements, `$lastN` returns all elements in the group.
// See https://www.mongodb.com/docs/manual/reference/operator/aggregation/lastN/
func LastN(input, n any) bson.Entry {
	return bson.E(operator.LastN, bson.D(
		bson.E("input", input),
		bson.E("n", n),
	))
}
