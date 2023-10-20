package executor

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	rawopts "go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"testing"
	"time"
)

const (
	testDBName   = "test_db"
	testCollName = "test_coll"
)

var (
	testCli  *mongo.Client
	testColl *mongo.Collection
)

type TestDoc struct {
	bson.Doc[primitive.ObjectID] `bson:"-"`
	Id                           primitive.ObjectID `bson:"_id"`
	Name                         string             `bson:"name"`
	Age                          int                `bson:"age"`
	CreateTime                   time.Time          `bson:"createTime"`
}

func (t *TestDoc) GetId() primitive.ObjectID {
	return t.Id
}

func init() {
	uri := os.Getenv("CI_MONGO_URI")
	var err error
	testCli, err = mongo.Connect(context.TODO(), rawopts.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	err = testCli.Database(testDBName).Collection(testCollName).Drop(context.TODO())
	if err != nil {
		panic(err)
	}
	testColl = testCli.Database(testDBName).Collection(testCollName)
}

func TestMain(m *testing.M) {
	_ = m.Run()
}
