package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"id"`
	Username string             `bson:"username"`
	Email    string             `bson:"email"`
}
