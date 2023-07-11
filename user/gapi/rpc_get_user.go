package gapi

import (
	"context"
	"errors"
	"firebase.google.com/go/v4/auth"
	"fmt"
	v1 "github.com/Dejan91/inventory-management/proto/user/api/v1"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetUser(ctx context.Context, r *v1.GetUserRequest) (*v1.GetUserResponse, error) {
	if violations := validateGetUserRequest(r); violations != nil {
		return nil, invalidArgumentError(violations)
	}

	user, err := s.db.GetUser(ctx, "_id", r.GetUid())
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, status.Errorf(codes.NotFound, "user not found: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to get user: %s", err)
	}

	firestoreUser, err := s.store.GetUser(ctx, user.ExternalID)
	if err != nil {
		if auth.IsUserNotFound(err) {
			return nil, status.Errorf(codes.NotFound, "user not found in firestore: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to get user from firestore: %s", err)
	}

	return &v1.GetUserResponse{
		Uid:         user.ID.Hex(),
		ExternalUid: user.ExternalID,
		Username:    user.Username,
		Email:       firestoreUser.Email,
	}, nil
}

func validateGetUserRequest(r *v1.GetUserRequest) []*errdetails.BadRequest_FieldViolation {
	var violations []*errdetails.BadRequest_FieldViolation

	if r.GetAuthUid() == "" {
		violations = append(violations, fieldViolation("auth_uid", fmt.Errorf("auth_uid can't be empty")))
	}

	if r.GetUid() == "" {
		violations = append(violations, fieldViolation("uid", fmt.Errorf("uid can't be empty")))
	}

	return violations
}
