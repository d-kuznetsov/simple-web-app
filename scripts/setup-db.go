package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	//"github.com/d-kuznetsov/chat/config"
)

type User struct {
	Username string
	Password string
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
		User{"Vanya", "vanya123"},
		User{"Katya", "katya123"},
	}
	insertResult, err := userCollection.InsertMany(ctx, users)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted Documents: ", insertResult)
}
