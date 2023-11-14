# Typed Mongo

Typed Mongo is an easy to use mongodb go driver, provided generic types. It's based on [mongo-go-driver](https://github.com/mongodb/mongo-go-driver).

## Requirements

- `Go 1.18` and higher
- `Mongo 3.6` and higher

## Installation

```shell
go get github.com/qianwj/typed/mongo
```

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
	"time"
)

type Person struct {
	model.Doc[primitive.ObjectID] `bson:"-"`
	Id                           primitive.ObjectID `bson:"_id,omitempty"`
	Name                         string             `bson:"name,omitempty"`
	Age                          int                `bson:"age,omitempty"`
}

func (p *Person) GetId() primitive.ObjectID {
	return p.Id
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
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

## Compare with `mongo-go-driver`

### Connected to server

1. mongo-go-driver
```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
if err != nil {
	log.Fatal(err)
}
err = client.Ping(ctx, readpref.Primary())
```

2. typed mongo
```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
cli, err := client.NewBuilder().ApplyUri("mongodb://localhost:27017").Ping(readpref.Primary()).Build(ctx)
```

### Insert Documents

1. mongo-go-driver
```go
ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
res, err := collection.InsertOne(ctx, bson.D{{"name", "pi"}, {"value", 3.14159}})
id := res.InsertedID
```

2. typed mongo
```go
ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
id, err := coll.InsertOne(&Person{Name: "alice", Age: 18}).Execute(ctx)
```

### Query Documents
