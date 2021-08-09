package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/d-kuznetsov/chat/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Id       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username,omitempty"`
	Password string             `bson:"password,omitempty"`
}

var client *mongo.Client

func Connect() *mongo.Client {
	var err error
	if client != nil {
		return client
	}
	clientOpts := options.Client().ApplyURI(config.MongoURI)
	client, err = mongo.NewClient(clientOpts)
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB is open.")
	return client
}

func Close() {
	if client == nil {
		return
	}

	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connection to MongoDB is closed.")
}

func FindUserByName(name string) (*User, error) {
	if client == nil {
		log.Fatal("There isn't db client")
	}
	var user User
	collection := client.Database("chat").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, bson.M{"username": name}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		fmt.Println("user does not exist")
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUser(username, password string) (*mongo.InsertOneResult, error) {
	if client == nil {
		log.Fatal("There isn't db client")
	}
	collection := client.Database("chat").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	user := User{
		Username: username,
		Password: password,
	}
	res, err := collection.InsertOne(ctx, user)
	fmt.Println("user was created")
	return res, err
}
