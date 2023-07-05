package store

import (
	"context"
	"firebase.google.com/go/v4/auth"
	"fmt"
)

type Store interface {
	CreateUser(context.Context, *auth.UserToCreate) (*auth.UserRecord, error)
	GetUser(ctx context.Context, uid string) (*auth.UserRecord, error)
	UpdateUser(ctx context.Context, uid string, user *auth.UserToUpdate) error
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
		return nil, fmt.Errorf("failed creating user in firestore: %w", err)
	}

	return userRecord, nil
}

func (f *FirebaseStore) GetUser(ctx context.Context, uid string) (*auth.UserRecord, error) {
	userRecord, err := f.client.GetUser(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("failed getting user from firestore: %w", err)
	}

	return userRecord, nil
}

func (f *FirebaseStore) UpdateUser(ctx context.Context, uid string, user *auth.UserToUpdate) error {
	_, err := f.client.UpdateUser(ctx, uid, user)
	if err != nil {
		return fmt.Errorf("failed updating user in firestore: %w", err)
	}

	return nil
}

func (f *FirebaseStore) DeleteUser(ctx context.Context, uid string) error {
	err := f.client.DeleteUser(ctx, uid)
	if err != nil {
		return fmt.Errorf("failed deleting user in firestore: %w", err)
	}

	return nil
}
