package expressions

import (
	"github.com/qianwj/typed/mongo/bson"
)

type Expression interface {
	~string | bson.Entry | bson.UnorderedMap | bson.Array | *bson.Map
}

type Number interface {
	bson.Number | Expression
}

func FieldPath(path string) string {
	return "$" + path
}
