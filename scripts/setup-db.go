package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/d-kuznetsov/chat/config"
	"github.com/d-kuznetsov/chat/models"
)

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI(config.MongoURI))
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
		models.User{Username: "jennifer", Password: "1234"},
		models.User{Username: "joey", Password: "1234"},
	}
	insertResult, err := userCollection.InsertMany(ctx, users)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted users: ", insertResult.InsertedIDs)
	articleCollection := client.Database("chat").Collection("articles")
	articles := []interface{}{
		models.Article{
			Title: "Golang",
			User:  insertResult.InsertedIDs[0].(string),
		},
		models.Article{
			Title: "JavaScript",
			User:  insertResult.InsertedIDs[1].(string),
		},
		models.Article{
			Title: "Python",
			User:  insertResult.InsertedIDs[1].(string),
		},
		models.Article{
			Title: "Rust",
			User:  insertResult.InsertedIDs[0].(string),
		},
		models.Article{
			Title: "Scala",
			User:  insertResult.InsertedIDs[0].(string),
		},
	}
	insertResult, err = articleCollection.InsertMany(ctx, articles)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted articles: ", insertResult.InsertedIDs)
}
