package service

import (
	"context"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"strings"
)

var noAuthRoutes = map[string]struct{}{
	"": {},
}

type AuthInterceptor struct {
	authClient *auth.Client
}

func NewAuthInterceptor(authClient *auth.Client) *AuthInterceptor {
	return &AuthInterceptor{
		authClient: authClient,
	}
}

func (a *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		log.Println("--> unary interceptor: ", info.FullMethod)

		if _, ok := noAuthRoutes[info.FullMethod]; !ok {
			newCtx, err := a.authenticate(ctx)
			if err != nil {
				return nil, err
			}

			return handler(newCtx, req)
		}

		return handler(ctx, req)
	}
}

func (a *AuthInterceptor) authenticate(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	splitToken := strings.Split(values[0], "Bearer ")
	accessToken := splitToken[1]
	token, err := a.authClient.VerifyIDToken(ctx, accessToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token is not valid: %s", err)
	}

	md.Append("auth_uid", token.UID)
	ctx = metadata.NewIncomingContext(ctx, md)

	return ctx, nil
}
