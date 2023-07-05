package gapi

import (
	"context"
	"firebase.google.com/go/v4/auth"
	"github.com/Dejan91/inventory-management/user/mocks"
	"github.com/Dejan91/inventory-management/user/model"
	"github.com/Dejan91/inventory-management/user/pb"
	"github.com/Dejan91/inventory-management/user/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestServer_CreateUser(t *testing.T) {
	ctx := context.Background()

	externalUID := util.RandomUID()
	uid := primitive.NewObjectID()
	email := util.RandomEmail()
	username := util.RandomUsername()
	password := util.RandomPassword()

	testCases := []struct {
		name          string
		r             *pb.CreateUserRequest
		buildStubs    func(db *mocks.MockDB, store *mocks.MockFirebaseStore)
		checkResponse func(t *testing.T, res *pb.CreateUserResponse, err error)
	}{
		{
			name: "OK",
			r: &pb.CreateUserRequest{
				AuthUid:  util.RandomUID(),
				Username: username,
				Email:    email,
				Password: password,
			},
			buildStubs: func(db *mocks.MockDB, store *mocks.MockFirebaseStore) {
				firestoreUser := &auth.UserRecord{
					UserInfo: &auth.UserInfo{
						UID:   externalUID,
						Email: email,
					},
				}

				db.On("GetUser", ctx, "username", username).
					Once().
					Return(nil, nil)

				store.On("CreateUser", ctx,
					(&auth.UserToCreate{}).
						Email(email).
						Password(password)).
					Once().
					Return(firestoreUser, nil)

				user := &model.DBUser{
					ExternalID: firestoreUser.UserInfo.UID,
					Username:   username,
				}
				db.On("CreateUser", ctx, user).
					Once().
					Return(&model.DBUser{
						ID:         uid,
						ExternalID: firestoreUser.UserInfo.UID,
						Username:   username,
					}, nil)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, res)

				assert.Equal(t, res.GetUsername(), username)
				assert.Equal(t, res.GetUid(), uid.Hex())
				assert.Equal(t, res.GetExternalUid(), externalUID)
				assert.Equal(t, res.GetEmail(), email)
			},
		},
		{
			"Duplicate username",
			&pb.CreateUserRequest{
				AuthUid:  util.RandomUID(),
				Username: username,
				Email:    email,
				Password: password,
			},
			func(db *mocks.MockDB, store *mocks.MockFirebaseStore) {
				firestoreUser := &auth.UserRecord{
					UserInfo: &auth.UserInfo{
						UID:   externalUID,
						Email: email,
					},
				}

				db.On("GetUser", ctx, "username", username).
					Once().
					Return(&model.DBUser{
						ID:         uid,
						ExternalID: firestoreUser.UserInfo.UID,
						Username:   username,
					}, nil)

				db.AssertNotCalled(t, "CreateUser")
				store.AssertNotCalled(t, "CreateUser")
			},
			func(t *testing.T, res *pb.CreateUserResponse, err error) {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.AlreadyExists, st.Code())
			},
		},
		{
			"Invalid email",
			&pb.CreateUserRequest{
				AuthUid:  util.RandomUID(),
				Username: username,
				Email:    "invalid-email",
				Password: password,
			},
			func(db *mocks.MockDB, store *mocks.MockFirebaseStore) {
				db.AssertNotCalled(t, "GetUser")
				db.AssertNotCalled(t, "CreateUser")
				store.AssertNotCalled(t, "CreateUser")
			},
			func(t *testing.T, res *pb.CreateUserResponse, err error) {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
		{
			"Internal error",
			&pb.CreateUserRequest{
				AuthUid:  util.RandomUID(),
				Username: username,
				Email:    email,
				Password: password,
			},
			func(db *mocks.MockDB, store *mocks.MockFirebaseStore) {
				db.On("GetUser", ctx, "username", username).
					Once().
					Return(nil, mongo.ErrClientDisconnected)

				db.AssertNotCalled(t, "CreateUser")
				store.AssertNotCalled(t, "CreateUser")
			},
			func(t *testing.T, res *pb.CreateUserResponse, err error) {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Internal, st.Code())
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		dbMock := new(mocks.MockDB)
		storeMock := new(mocks.MockFirebaseStore)

		tc.buildStubs(dbMock, storeMock)
		server := newTestServer(dbMock, storeMock)
		response, err := server.CreateUser(ctx, tc.r)

		tc.checkResponse(t, response, err)
	}
}
