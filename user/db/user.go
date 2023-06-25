package db

import (
	"context"
	"fmt"
	"github.com/Dejan91/inventory-management/user/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (m *MongoDB) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	m.client.Database(m.dbName).Collection("users")
	coll := m.client.Database(m.dbName).Collection("users")
	filter := bson.D{{"email", email}}

	var result model.User
	err := coll.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no results %w", err)
		}

		return nil, fmt.Errorf("error while getting user from database %w", err)
	}

	return &result, nil
}

func (m *MongoDB) Create(ctx context.Context, user *model.User) (*model.User, error) {
	createdUser, err := m.client.Database(m.dbName).Collection("users").InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	if oid, ok := createdUser.InsertedID.(primitive.ObjectID); ok {
		user.ID = oid
	}

	return user, nil
}
