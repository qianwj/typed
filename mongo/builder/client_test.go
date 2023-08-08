package builder

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"testing"
)

func TestClientBuilder_ApplyUri(t *testing.T) {
	cli, err := NewClient().ApplyUri("mongodb://localhost:27017").Build(context.TODO())
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if err := cli.Ping(context.TODO(), readpref.Primary()); err != nil {
		t.Error(err)
		t.FailNow()
	}
}
