package gapi

import (
	"context"
	"errors"
	"firebase.google.com/go/v4/auth"
	"fmt"
	v1 "github.com/Dejan91/inventory-management/proto/user/api/v1"
	"github.com/Dejan91/inventory-management/user/model"
	"github.com/Dejan91/inventory-management/user/val"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateUser(ctx context.Context, r *v1.UpdateUserRequest) (*v1.UpdateUserResponse, error) {
	if violations := validateUpdateUserRequest(r); violations != nil {
		return nil, invalidArgumentError(violations)
	}

	primitiveUID, err := primitive.ObjectIDFromHex(r.GetUid())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to udpate user: %s", err)
	}

	u, err := s.db.GetUser(ctx, "_id", r.GetUid())
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, status.Errorf(codes.NotFound, "user not found in db: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to udpate user: %s", err)
	}

	firestoreUser, err := s.store.GetUser(ctx, u.ExternalID)
	if err != nil {
		if auth.IsUserNotFound(err) {
			return nil, status.Errorf(codes.NotFound, "user not found in firestore: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to udpate user: %s", err)
	}

	err = s.db.Transaction(ctx, func() error {
		err = s.db.UpdateUser(ctx, &model.DBUser{
			ID:       primitiveUID,
			Username: r.GetUsername(),
		})
		if err != nil {
			return err
		}

		if firestoreUser.Email != r.GetEmail() {
			err = s.store.UpdateUser(ctx, u.ExternalID, (&auth.UserToUpdate{}).Email(r.GetEmail()))
		}

		return nil
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	return &v1.UpdateUserResponse{
		Uid:         u.ID.Hex(),
		ExternalUid: u.ExternalID,
		Username:    r.GetUsername(),
		Email:       r.GetEmail(),
	}, nil
}

func validateUpdateUserRequest(r *v1.UpdateUserRequest) []*errdetails.BadRequest_FieldViolation {
	var violations []*errdetails.BadRequest_FieldViolation

	if r.GetAuthUid() == "" {
		violations = append(violations, fieldViolation("auth_uid", fmt.Errorf("auth_uid can't be empty")))
	}

	if err := val.ValidateEmail(r.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	if err := val.ValidateUsername(r.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	return violations
}
