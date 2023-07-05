package gapi

import (
	"context"
	"github.com/Dejan91/inventory-management/user/mocks"
	"github.com/Dejan91/inventory-management/user/model"
	"github.com/Dejan91/inventory-management/user/pb"
	"github.com/Dejan91/inventory-management/user/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestServer_DeleteUser(t *testing.T) {
	ctx := context.Background()

	uid := primitive.NewObjectID()
	externalUID := util.RandomUID()
	username := util.RandomUsername()

	testCases := []struct {
		name          string
		r             *pb.DeleteUserRequest
		buildStubs    func(db *mocks.MockDB, store *mocks.MockFirebaseStore)
		checkResponse func(t *testing.T, res *pb.DeleteUserResponse, err error)
	}{
		{
			name: "OK",
			r: &pb.DeleteUserRequest{
				AuthUid: util.RandomUID(),
				Uid:     uid.Hex(),
			},
			buildStubs: func(db *mocks.MockDB, store *mocks.MockFirebaseStore) {
				db.On("GetUser", ctx, "_id", uid.Hex()).
					Once().
					Return(&model.DBUser{
						ID:         uid,
						ExternalID: externalUID,
						Username:   username,
					}, nil)

				db.On("Transaction", ctx, mock.AnythingOfType("func() error")).
					Once().
					Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteUserResponse, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, res)

				assert.Equal(t, res.Message, returnMessage)
			},
		},
		{
			name: "Invalid arguments",
			r: &pb.DeleteUserRequest{
				AuthUid: util.RandomUID(),
				Uid:     "",
			},
			buildStubs: func(db *mocks.MockDB, store *mocks.MockFirebaseStore) {
				db.AssertNotCalled(t, "GetUser")
				db.AssertNotCalled(t, "Transaction")
			},
			checkResponse: func(t *testing.T, res *pb.DeleteUserResponse, err error) {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
		{
			name: "Not found",
			r: &pb.DeleteUserRequest{
				AuthUid: util.RandomUID(),
				Uid:     uid.Hex(),
			},
			buildStubs: func(db *mocks.MockDB, store *mocks.MockFirebaseStore) {
				db.On("GetUser", ctx, "_id", uid.Hex()).
					Once().
					Return(nil, mongo.ErrNoDocuments)

				db.AssertNotCalled(t, "Transaction")
			},
			checkResponse: func(t *testing.T, res *pb.DeleteUserResponse, err error) {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, codes.NotFound, st.Code())
			},
		},
		{
			name: "Internal error",
			r: &pb.DeleteUserRequest{
				AuthUid: util.RandomUID(),
				Uid:     uid.Hex(),
			},
			buildStubs: func(db *mocks.MockDB, store *mocks.MockFirebaseStore) {
				db.On("GetUser", ctx, "_id", uid.Hex()).
					Once().
					Return(nil, mongo.ErrClientDisconnected)

				db.AssertNotCalled(t, "Transaction")
			},
			checkResponse: func(t *testing.T, res *pb.DeleteUserResponse, err error) {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, codes.Internal, st.Code())
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		dbMock := new(mocks.MockDB)
		storeMock := new(mocks.MockFirebaseStore)

		tc.buildStubs(dbMock, storeMock)
		server := newTestServer(dbMock, storeMock)
		response, err := server.DeleteUser(ctx, tc.r)

		tc.checkResponse(t, response, err)
	}
}
