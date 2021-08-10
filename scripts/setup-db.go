package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	//"github.com/d-kuznetsov/chat/config"
)

type User struct {
	Id       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username,omitempty"`
	Password string             `bson:"password,omitempty"`
}

type Article struct {
	Id    primitive.ObjectID `bson:"_id,omitempty"`
	Title string             `bson:"title,omitempty"`
	Date  string             `bson:"date,omitempty"`
	Text  string             `bson:"text,omitempty"`
	User  primitive.ObjectID `bson:"user,omitempty"`
}

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	err = client.Database("chat").Drop(ctx)
	if err != nil {
		log.Fatal(err)
	}
	userCollection := client.Database("chat").Collection("users")
	users := []interface{}{
		User{Username: "vanya", Password: "1234"},
		User{Username: "katya", Password: "1234"},
	}
	insertResult, err := userCollection.InsertMany(ctx, users)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted users: ", insertResult.InsertedIDs)
	articleCollection := client.Database("chat").Collection("articles")
	articles := []interface{}{
		Article{
			Title: "Golang",
			User:  insertResult.InsertedIDs[0].(primitive.ObjectID),
		},
		Article{
			Title: "JavaScript",
			User:  insertResult.InsertedIDs[1].(primitive.ObjectID),
		},
		Article{
			Title: "Python",
			User:  insertResult.InsertedIDs[1].(primitive.ObjectID),
		},
		Article{
			Title: "Rust",
			User:  insertResult.InsertedIDs[0].(primitive.ObjectID),
		},
		Article{
			Title: "Scala",
			User:  insertResult.InsertedIDs[1].(primitive.ObjectID),
		},
	}
	insertResult, err = articleCollection.InsertMany(ctx, articles)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted articles: ", insertResult.InsertedIDs)
}
