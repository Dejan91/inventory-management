package server

import (
	"context"
	pb "github.com/Dejan91/inventory-management/proto/gateway/user/api/v1"
	userV1 "github.com/Dejan91/inventory-management/proto/user/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServer struct {
	userClient userV1.UserClient
	pb.UnimplementedUserServer
}

func NewUserServer(userClient userV1.UserClient) *UserServer {
	return &UserServer{
		userClient: userClient,
	}
}

func (u *UserServer) CreateUser(ctx context.Context, r *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	authUID, err := getAuthUID(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "couldn't get auth id from request, :%s", err)
	}

	res, err := u.userClient.CreateUser(ctx, &userV1.CreateUserRequest{
		AuthUid:  authUID,
		Username: r.GetUsername(),
		Email:    r.GetEmail(),
		Password: r.GetPassword(),
	})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return nil, status.Errorf(codes.Internal, "failed to create user")
		}
		return nil, status.Errorf(st.Code(), "failed to create user: %st", err)
	}

	return &pb.CreateUserResponse{
		Uid:      res.GetUid(),
		Username: res.GetUsername(),
		Email:    res.GetEmail(),
	}, nil
}

func (u *UserServer) GetUser(ctx context.Context, r *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	authUID, err := getAuthUID(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "couldn't get auth id from request, :%s", err)
	}

	res, err := u.userClient.GetUser(ctx, &userV1.GetUserRequest{
		AuthUid: authUID,
		Uid:     r.GetUid(),
	})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return nil, status.Errorf(codes.Internal, "failed to retreive user")
		}
		return nil, status.Errorf(st.Code(), "failed to retreive user: %s", err)
	}

	return &pb.GetUserResponse{
		Uid:         res.GetUid(),
		ExternalUid: res.GetExternalUid(),
		Username:    res.GetUsername(),
		Email:       res.GetEmail(),
	}, nil
}

func (u *UserServer) UpdateUser(ctx context.Context, r *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	authUID, err := getAuthUID(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "couldn't get auth id from request, :%s", err)
	}

	res, err := u.userClient.UpdateUser(ctx, &userV1.UpdateUserRequest{
		AuthUid:  authUID,
		Uid:      r.GetUid(),
		Username: r.GetUsername(),
		Email:    r.GetEmail(),
	})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return nil, status.Errorf(codes.Internal, "failed to update user")
		}
		return nil, status.Errorf(st.Code(), "failed to update user: %s", err)
	}

	return &pb.UpdateUserResponse{
		Uid:      res.GetUid(),
		Username: res.GetUsername(),
		Email:    res.GetEmail(),
	}, nil
}

func (u *UserServer) DeleteUser(ctx context.Context, r *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	authUID, err := getAuthUID(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "couldn't get auth id from request, :%s", err)
	}

	res, err := u.userClient.DeleteUser(ctx, &userV1.DeleteUserRequest{
		AuthUid: authUID,
		Uid:     r.GetUid(),
	})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return nil, status.Errorf(codes.Internal, "failed to delete user")
		}
		return nil, status.Errorf(st.Code(), "failed to delete user: %s", err)
	}

	return &pb.DeleteUserResponse{
		Message: res.Message,
	}, nil
}
