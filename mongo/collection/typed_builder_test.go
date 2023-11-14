package collection

import (
	"github.com/qianwj/typed/mongo/model"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"testing"
)

type mockDoc struct {
	model.Doc[primitive.ObjectID] `bson:"-"`
	ID                            primitive.ObjectID `bson:"_id"`
	Name                          string             `bson:"name"`
}

func (m *mockDoc) GetId() primitive.ObjectID {
	return m.ID
}

func TestTypedBuilder_ReadPreference(t *testing.T) {
	builder := NewTypedBuilder[*mockDoc, primitive.ObjectID](nil, "test").ReadPreference(readpref.Secondary())
	actual := builder.opts
	expected := options.Collection().SetReadPreference(readpref.Secondary())
	assert.Equal(t, expected, actual)
}
