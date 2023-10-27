# Typed Mongo

An easy to use mongodb go driver. provided generic types.

## Quick Start

```go
package main

import (
	"context"
	"github.com/qianwj/typed/mongo/bson"
	"github.com/qianwj/typed/mongo/client"
	"github.com/qianwj/typed/mongo/collection"
	"github.com/qianwj/typed/mongo/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

type Person struct {
	bson.Doc[primitive.ObjectID] `bson:"-"`
	Id   primitive.ObjectID      `bson:"_id,omitempty"`
	Name string                  `bson:"name,omitempty"`
	Age  int                     `bson:"age,omitempty"`
}

func (p *Person) GetId() primitive.ObjectID {
	return p.Id
}

func main() {
	ctx := context.TODO()
	cli, err := client.NewBuilder().ApplyUri("mongodb://localhost:27017/admin").Build(ctx)
	if err != nil {
		log.Fatal(err)
	}
	db := cli.DefaultDatabase().Build().Raw()
	coll := collection.NewTypedBuilder[*Person, primitive.ObjectID](db, "person").Build()
	id, err := coll.InsertOne(&Person{Name: "alice", Age: 18}).Execute(ctx)
	if err != nil {
		log.Fatal(err)
	}
	dbPerson, err := coll.FindOneById(id).Execute(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("person: %+v", dbPerson)
}
```