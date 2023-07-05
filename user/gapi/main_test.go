package gapi

import (
	"github.com/Dejan91/inventory-management/user/mocks"
	"github.com/Dejan91/inventory-management/user/util"
)

func newTestServer(dbMock *mocks.MockDB, firebaseMock *mocks.MockFirebaseStore) *Server {
	config := &util.Config{}

	return NewServer(config, dbMock, firebaseMock)
}
