package db

import (
	"context"
	"fmt"
	"github.com/Dejan91/inventory-management/user/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (m *MongoDB) CreateUser(ctx context.Context, user *model.DBUser) (*model.DBUser, error) {
	createdUser, err := m.client.Database(m.dbName).Collection("users").InsertOne(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user in mongodb: %w", err)
	}

	if oid, ok := createdUser.InsertedID.(primitive.ObjectID); ok {
		user.ID = oid
	}

	return user, nil
}

func (m *MongoDB) GetUser(ctx context.Context, field, value string) (*model.DBUser, error) {
	var filter = bson.D{}
	if field == "_id" {
		v, err := primitive.ObjectIDFromHex(value)
		if err != nil {
			return nil, err
		}
		filter = bson.D{{field, v}}
	} else {
		filter = bson.D{{field, value}}
	}

	coll := m.client.Database(m.dbName).Collection("users")

	var result model.DBUser
	err := coll.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no results %w", err)
		}

		return nil, fmt.Errorf("error while getting user from mongodb: %w", err)
	}

	return &result, nil
}

func (m *MongoDB) UpdateUser(ctx context.Context, user *model.DBUser) error {
	filter := bson.M{"_id": user.ID}
	update := bson.D{{"$set", bson.D{
		{"username", user.Username},
	}}}

	_, err := m.client.Database(m.dbName).Collection("users").UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed updating user in mongodb: %w", err)
	}

	return nil
}

func (m *MongoDB) DeleteUser(ctx context.Context, uid string) error {
	idPrimitive, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		return fmt.Errorf("failed deleting user in mongodb: %w", err)
	}

	filter := bson.M{"_id": idPrimitive}
	_, err = m.client.Database(m.dbName).Collection("users").DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed deleting user in mongodb: %w", err)
	}

	return nil
}
