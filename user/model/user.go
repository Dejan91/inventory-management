package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type DBUser struct {
	ID         primitive.ObjectID `bson:"_id"`
	ExternalID string             `bson:"external_id"`
	Username   string             `bson:"username"`
}
