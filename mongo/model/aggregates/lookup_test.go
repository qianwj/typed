package aggregates

import (
	"github.com/qianwj/typed/mongo/bson"
	"github.com/qianwj/typed/mongo/model/filters"
	"github.com/qianwj/typed/mongo/operator"
	"testing"
)

func TestNewLookup(t *testing.T) {
	l := NewLookup("holidays", "holidays").Pipeline(bson.A{
		Match(filters.Eq("year", 2018)).Marshal(),
		Project(bson.M{"_id": 0, "date": bson.M{"name": "$name", "date": "$date"}}).Marshal(),
		bson.M{operator.ReplaceRoot: bson.M{"newRoot": "$date"}},
	})
	t.Logf("parse: %+v", l.ToMap())
}
