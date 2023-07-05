package mocks

import (
	"context"
	"github.com/Dejan91/inventory-management/user/model"
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Transaction(ctx context.Context, fn func() error) error {
	args := m.Called(ctx, fn)

	var r0 error
	if args.Get(0) != nil {
		r0 = args.Get(0).(error)
	}

	return r0
}

func (m *MockDB) Close(ctx context.Context) error {
	args := m.Called(ctx)

	var r0 error
	if args.Get(0) != nil {
		r0 = args.Get(0).(error)
	}

	return r0
}

func (m *MockDB) CreateUser(ctx context.Context, user *model.DBUser) (*model.DBUser, error) {
	args := m.Called(ctx, user)

	var r0 *model.DBUser
	if args.Get(0) != nil {
		r0 = args.Get(0).(*model.DBUser)
	}

	var r1 error
	if args.Get(1) != nil {
		r1 = args.Get(1).(error)
	}

	return r0, r1
}

func (m *MockDB) GetUser(ctx context.Context, field, value string) (*model.DBUser, error) {
	args := m.Called(ctx, field, value)

	var r0 *model.DBUser
	if args.Get(0) != nil {
		r0 = args.Get(0).(*model.DBUser)
	}

	var r1 error
	if args.Get(1) != nil {
		r1 = args.Get(1).(error)
	}

	return r0, r1
}

func (m *MockDB) UpdateUser(ctx context.Context, user *model.DBUser) error {
	args := m.Called(ctx, user)

	var r0 error
	if args.Get(0) != nil {
		r0 = args.Get(0).(error)
	}

	return r0
}

func (m *MockDB) DeleteUser(ctx context.Context, uid string) error {
	args := m.Called(ctx, uid)

	var r0 error
	if args.Get(1) != nil {
		r0 = args.Get(1).(error)
	}

	return r0
}
