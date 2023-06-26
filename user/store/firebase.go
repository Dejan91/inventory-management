package store

import (
	"context"
	"firebase.google.com/go/v4/auth"
)

type Store interface {
	CreateUser(context.Context, *auth.UserToCreate) (*auth.UserRecord, error)
}

type FirebaseStore struct {
	client *auth.Client
}

func NewFirebaseStore(client *auth.Client) Store {
	return &FirebaseStore{
		client: client,
	}
}

func (f *FirebaseStore) CreateUser(ctx context.Context, user *auth.UserToCreate) (*auth.UserRecord, error) {
	userRecord, err := f.client.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return userRecord, nil
}
