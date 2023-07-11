package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
)

func getAuthUID(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("error getting metadata from incomint context")
	}

	authUID := md["auth_uid"]
	return authUID[0], nil
}
