package aggregate

import (
	"github.com/qianwj/typed/mongo/model/operator"
	"go.mongodb.org/mongo-driver/bson"
)

func Abs(root bson.M, expression bson.M) {
	root[operator.Abs] = expression
}
