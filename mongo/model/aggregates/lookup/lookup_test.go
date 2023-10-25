package lookup

import (
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestNewLookup(t *testing.T) {
	expected := New("holidays", "holidays")
	bytes, err := bson.Marshal(expected)
	if err != nil {
		t.Errorf("marshal lookup error: %+v", err)
		t.FailNow()
	}
	var actual JoinCondition
	if err := bson.Unmarshal(bytes, &actual); err != nil {
		t.Errorf("unmarshal lookup error: %+v", err)
	}
	assert.Equal(t, expected, &actual)
}
