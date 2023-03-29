package database

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var Ctx = context.TODO()

var PostsCollection *mongo.Collection

func ConnectDB() {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		panic("You must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	client, err := mongo.Connect(Ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	if err := client.Ping(Ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	db := client.Database("syo")
	PostsCollection = db.Collection("posts")

	fmt.Println("Database connected")
}
