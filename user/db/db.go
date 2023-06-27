package db

import (
	"context"
	"fmt"
	"github.com/Dejan91/inventory-management/user/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DB interface {
	Close(ctx context.Context) error
	GetUserByEmail(ctx context.Context, email string) (*model.DBUser, error)
	GetUserByUsername(ctx context.Context, username string) (*model.DBUser, error)
	CreateUser(ctx context.Context, user *model.DBUser) (*model.DBUser, error)
	DeleteUser(ctx context.Context, uid string) error
}

type MongoDB struct {
	client *mongo.Client
	dbName string
}

func NewMongoDB(username, password, host, dbName string) (DB, error) {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:27017/?maxPoolSize=20&w=majority",
		username,
		password,
		host,
	)

	ctx := context.Background()
	fmt.Println("connecting to mongodb")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	fmt.Println("pinging mongodb")
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("connected to mongodb")
	return &MongoDB{
		client: client,
		dbName: dbName,
	}, nil
}

func (m *MongoDB) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}
