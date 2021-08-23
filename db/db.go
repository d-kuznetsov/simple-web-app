package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/d-kuznetsov/blog/config"
)

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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
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

func checkClient() {
	if client == nil {
		log.Fatal("There isn't db client")
	}
}

func FindUserByName(name string) (*User, error) {
	checkClient()
	var user User
	collection := client.Database("blog").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, bson.M{"username": name}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		fmt.Println("User does not exist")
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUser(username, password string) (primitive.ObjectID, error) {
	checkClient()
	collection := client.Database("blog").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	user := User{
		Username: username,
		Password: password,
	}
	res, err := collection.InsertOne(ctx, user)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), err
}

func GetAllArticles() ([]*Article, error) {
	checkClient()
	collection := client.Database("blog").Collection("articles")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, bson.M{})
	if err == mongo.ErrNoDocuments {
		fmt.Println("There are not any articles")
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	var articles []*Article
	err = cursor.All(ctx, &articles)
	if err != nil {
		return nil, err
	}
	return articles, nil
}

func GetArticleById(id string) (*Article, error) {
	checkClient()
	collection := client.Database("blog").Collection("articles")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var article Article
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	err = collection.FindOne(ctx, bson.M{"_id": objId}).Decode(&article)
	if err == mongo.ErrNoDocuments {
		fmt.Println("Article does not exist")
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &article, nil
}

func CreateArticle(title, text, userId string) (primitive.ObjectID, error) {
	checkClient()
	collection := client.Database("blog").Collection("articles")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(userId)
	article := Article{
		Title: title,
		Text:  text,
		User:  objId,
	}
	res, err := collection.InsertOne(ctx, article)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), err
}

func UpdateArticle(id, title, text string) (*mongo.UpdateResult, error) {
	checkClient()
	collection := client.Database("blog").Collection("articles")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(id)
	res, err := collection.UpdateByID(ctx, objId, bson.D{
		{"$set", bson.M{"title": title, "text": text}},
	})
	if err == mongo.ErrNoDocuments {
		fmt.Println("There are not any articles")
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return res, err
}

func GetArticlesByUserId(userId string) ([]*Article, error) {
	checkClient()
	collection := client.Database("blog").Collection("articles")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(userId)
	cursor, err := collection.Find(ctx, bson.M{"user": objId})
	if err == mongo.ErrNoDocuments {
		fmt.Println("There are not any articles")
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	var articles []*Article
	err = cursor.All(ctx, &articles)
	if err != nil {
		return nil, err
	}
	return articles, nil
}

func DeleteArticlesByIds(ids []string) (*mongo.DeleteResult, error) {
	checkClient()
	collection := client.Database("blog").Collection("articles")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var objIds []primitive.ObjectID

	for _, val := range ids {
		objId, _ := primitive.ObjectIDFromHex(val)
		objIds = append(objIds, objId)
	}

	res, err := collection.DeleteMany(ctx, bson.M{
		"_id": bson.M{
			"$in": objIds,
		},
	})
	if err == mongo.ErrNoDocuments {
		fmt.Println("There are not any articles")
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return res, err
}
