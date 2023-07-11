package db

import (
	"context"
	"fmt"
	"github.com/Dejan91/inventory-management/user/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type DB interface {
	CreateUser(ctx context.Context, user *model.DBUser) (*model.DBUser, error)
	GetUser(ctx context.Context, field, value string) (*model.DBUser, error)
	UpdateUser(ctx context.Context, user *model.DBUser) error
	DeleteUser(ctx context.Context, uid string) error
	Transaction(ctx context.Context, fn func() error) error
	Close(ctx context.Context) error
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
	// TODO:: Do once
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

func (m *MongoDB) Transaction(ctx context.Context, fn func() error) error {
	wc := writeconcern.Majority()
	txnOptions := options.Transaction().SetWriteConcern(wc)

	session, err := m.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	err = mongo.WithSession(ctx, session, func(ctx mongo.SessionContext) error {
		if err = session.StartTransaction(txnOptions); err != nil {
			return err
		}

		err = fn()
		if err != nil {
			return err
		}

		if err = session.CommitTransaction(ctx); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		if err = session.AbortTransaction(ctx); err != nil {
			return err
		}
		return err
	}

	return nil
}

func (m *MongoDB) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}
