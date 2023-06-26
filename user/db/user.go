package db

import (
	"context"
	"fmt"
	"github.com/Dejan91/inventory-management/user/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (m *MongoDB) GetUser(ctx context.Context, uid string) (*model.DBUser, error) {
	m.client.Database(m.dbName).Collection("users")
	coll := m.client.Database(m.dbName).Collection("users")
	filter := bson.D{{"external_id", uid}}

	var result model.DBUser
	err := coll.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no results %w", err)
		}

		return nil, fmt.Errorf("error while getting user from database %w", err)
	}

	return &result, nil
}

func (m *MongoDB) CreateUser(ctx context.Context, user *model.DBUser) (*model.DBUser, error) {
	createdUser, err := m.client.Database(m.dbName).Collection("users").InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	if oid, ok := createdUser.InsertedID.(primitive.ObjectID); ok {
		user.ID = oid
	}

	return user, nil
}
