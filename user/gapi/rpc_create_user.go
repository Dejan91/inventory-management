package gapi

import (
	"context"
	"errors"
	"firebase.google.com/go/v4/auth"
	"fmt"
	v1 "github.com/Dejan91/inventory-management/proto/user/api/v1"
	"github.com/Dejan91/inventory-management/user/model"
	"github.com/Dejan91/inventory-management/user/val"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateUser(ctx context.Context, r *v1.CreateUserRequest) (*v1.CreateUserResponse, error) {
	if violations := validateCreateUserRequest(r); violations != nil {
		return nil, invalidArgumentError(violations)
	}

	u, err := s.db.GetUser(ctx, "username", r.GetUsername())
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
		}
	}
	if u != nil {
		return nil, status.Errorf(codes.AlreadyExists, "user with provided username already exists")
	}

	userToCreate := (&auth.UserToCreate{}).
		Email(r.GetEmail()).
		Password(r.GetPassword())

	firebaseUser, err := s.store.CreateUser(ctx, userToCreate)
	if err != nil {
		if auth.IsEmailAlreadyExists(err) {
			return nil, status.Errorf(codes.AlreadyExists, "failed to create user: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	user := model.DBUser{
		ExternalID: firebaseUser.UID,
		Username:   r.GetUsername(),
	}

	dbUser, err := s.db.CreateUser(ctx, &user)
	if err != nil {
		if err = s.store.DeleteUser(ctx, firebaseUser.UID); err != nil {
			return nil, status.Errorf(
				codes.Internal,
				"error encountered failed to reverse created firebase user: %s", err,
			)
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	response := &v1.CreateUserResponse{
		Uid:         dbUser.ID.Hex(),
		ExternalUid: firebaseUser.UID,
		Username:    dbUser.Username,
		Email:       firebaseUser.Email,
	}

	return response, nil
}

func validateCreateUserRequest(r *v1.CreateUserRequest) []*errdetails.BadRequest_FieldViolation {
	var violations []*errdetails.BadRequest_FieldViolation

	if r.GetAuthUid() == "" {
		violations = append(violations, fieldViolation("uid", fmt.Errorf("invalid uid")))
	}

	if err := val.ValidateUsername(r.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := val.ValidateEmail(r.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	if err := val.ValidatePassword(r.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	return violations
}
