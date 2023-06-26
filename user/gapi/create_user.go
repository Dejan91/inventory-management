package gapi

import (
	"context"
	"firebase.google.com/go/v4/auth"
	"fmt"
	pb "github.com/Dejan91/inventory-management/proto/golang/user"
	"github.com/Dejan91/inventory-management/user/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"time"
)

func (s *Server) CreateUser(ctx context.Context, r *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	userToCreate := (&auth.UserToCreate{}).
		Email(r.GetEmail()).
		Password(r.GetPassword())

	firebaseUser, err := s.store.CreateUser(ctx, userToCreate)
	if err != nil {
		log.Fatalf("error creating user: %v\n", err)
	}

	fmt.Printf("%+v\n", firebaseUser)
	user := model.DBUser{
		Username: r.GetUsername(),
	}

	createdUser, err := s.db.CreateUser(ctx, &user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	response := &pb.CreateUserResponse{
		Username:  createdUser.Username,
		Email:     createdUser.Email,
		Password:  "password",
		CreatedAt: timestamppb.New(time.Now()),
	}

	return response, nil
}
