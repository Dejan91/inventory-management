package mocks

import (
	"context"
	"firebase.google.com/go/v4/auth"
	"github.com/stretchr/testify/mock"
)

type MockFirebaseStore struct {
	mock.Mock
}

func (m *MockFirebaseStore) CreateUser(ctx context.Context, user *auth.UserToCreate) (*auth.UserRecord, error) {
	args := m.Called(ctx, user)

	var r0 *auth.UserRecord
	if args.Get(0) != nil {
		r0 = args.Get(0).(*auth.UserRecord)
	}

	var r1 error
	if args.Get(1) != nil {
		r1 = args.Get(1).(error)
	}

	return r0, r1
}

func (m *MockFirebaseStore) GetUser(ctx context.Context, uid string) (*auth.UserRecord, error) {
	args := m.Called(ctx, uid)

	return args.Get(0).(*auth.UserRecord), args.Get(1).(error)
}

func (m *MockFirebaseStore) UpdateUser(ctx context.Context, uid string, user *auth.UserToUpdate) error {
	args := m.Called(ctx, uid, user)

	return args.Get(0).(error)
}

func (m *MockFirebaseStore) DeleteUser(ctx context.Context, uid string) error {
	args := m.Called(ctx, uid)

	var r0 error
	if args.Get(0) != nil {
		r0 = args.Get(0).(error)
	}

	return r0
}
