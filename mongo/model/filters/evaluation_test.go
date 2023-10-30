package filters

import (
	"github.com/qianwj/typed/mongo/operator"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestLike(t *testing.T) {
	expected, _ := bson.Marshal(bson.D{{Key: "a", Value: bson.D{{Key: operator.Regex, Value: primitive.Regex{
		Pattern: "abc",
		Options: "im",
	}}}}})
	actual, _ := bson.Marshal(Like("a", NewMatcher("abc").IgnoreCase().MultilineMatch()))
	assert.Equal(t, expected, actual)
}
