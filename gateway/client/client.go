package client

import (
	"context"
	"fmt"
	userV1 "github.com/Dejan91/inventory-management/proto/user/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	userSrv = "user:9090"
)

func NewUserClient(ctx context.Context) (userV1.UserClient, error) {
	conn, err := grpc.DialContext(
		ctx, userSrv, grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("user client failed: %w", err)
	}

	return userV1.NewUserClient(conn), nil
}
