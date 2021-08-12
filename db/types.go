package db

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
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
