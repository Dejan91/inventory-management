package db

import (
	"context"
	"fmt"
	"github.com/Dejan91/inventory-management/user/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (m *MongoDB) getUser(ctx context.Context, field, value string) (*model.DBUser, error) {
	filter := bson.D{{field, value}}
	coll := m.client.Database(m.dbName).Collection("users")

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

func (m *MongoDB) GetUserByEmail(ctx context.Context, email string) (*model.DBUser, error) {
	user, err := m.getUser(ctx, "email", email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (m *MongoDB) GetUserByUsername(ctx context.Context, username string) (*model.DBUser, error) {
	user, err := m.getUser(ctx, "username", username)
	if err != nil {
		return nil, err
	}

	return user, nil
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

func (m *MongoDB) DeleteUser(ctx context.Context, uid string) error {
	filter := bson.D{{"uid", uid}}
	_, err := m.client.Database(m.dbName).Collection("users").DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
