package executor

import (
	"context"
	"github.com/qianwj/typed/mongo/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
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
	err := mtest.Setup(mtest.NewSetupOptions())
	if err != nil {
		panic(err)
	}
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}
	cli, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	err = cli.Database(testDBName).Drop(context.TODO())
	if err != nil {
		panic(err)
	}
	testCli, testColl = cli, cli.Database(testDBName).Collection(testCollName)
}
