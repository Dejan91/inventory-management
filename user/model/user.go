package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type DBUser struct {
	ID         primitive.ObjectID `bson:"id"`
	ExternalID string             `bson:"external_id"`
	Username   string             `bson:"username"`
}

type User struct {
	DBUser
	Email string
}
