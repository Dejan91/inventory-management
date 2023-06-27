package store

import (
	"context"
	"firebase.google.com/go/v4/auth"
)

type Store interface {
	CreateUser(context.Context, *auth.UserToCreate) (*auth.UserRecord, error)
	DeleteUser(ctx context.Context, uid string) error
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

func (f *FirebaseStore) DeleteUser(ctx context.Context, uid string) error {
	return f.client.DeleteUser(ctx, uid)
}
