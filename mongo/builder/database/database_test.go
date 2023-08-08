package database

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"testing"
)

const testMongoUri = "mongodb://localhost:27017"

func TestDatabaseBuilder_ReadConcern(t *testing.T) {
	cli, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testMongoUri))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	db := NewDatabase(cli, "test_db").ReadConcern(readconcern.Majority()).Build()
	assert.Equal(t, readconcern.Majority(), db.Raw().ReadConcern())
}

func TestDatabaseBuilder_ReadPreference(t *testing.T) {
	cli, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testMongoUri))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	db := NewDatabase(cli, "test_db").ReadPreference(readpref.Primary()).Build()
	assert.Equal(t, readpref.Primary(), db.Raw().ReadPreference())
}

func TestDatabaseBuilder_Build(t *testing.T) {
	cli, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testMongoUri))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	db := NewDatabase(cli, "test_db").Build()
	collectionNames, err := db.Raw().ListCollectionNames(context.TODO(), bson.M{})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("collection names: %+v", collectionNames)
}
