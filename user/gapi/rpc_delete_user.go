package gapi

import (
	"context"
	"errors"
	"fmt"
	"github.com/Dejan91/inventory-management/user/pb"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const returnMessage = "User deleted"

func (s *Server) DeleteUser(ctx context.Context, r *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	if violations := validateDeleteUserRequest(r); violations != nil {
		return nil, invalidArgumentError(violations)
	}

	u, err := s.db.GetUser(ctx, "_id", r.GetUid())
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete user: %s", err)
	}

	err = s.db.Transaction(ctx, func() error {
		err = s.db.DeleteUser(ctx, r.GetUid())
		if err != nil {
			return err
		}

		err = s.store.DeleteUser(ctx, u.ExternalID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user: %s", err)
	}

	return &pb.DeleteUserResponse{
		Message: returnMessage,
	}, nil
}

func validateDeleteUserRequest(r *pb.DeleteUserRequest) []*errdetails.BadRequest_FieldViolation {
	var violations []*errdetails.BadRequest_FieldViolation

	if r.GetAuthUid() == "" {
		violations = append(violations, fieldViolation("auth_uid", fmt.Errorf("auth_uid can't be empty")))
	}

	if r.GetUid() == "" {
		violations = append(violations, fieldViolation("uid", fmt.Errorf("uid can't be empty")))
	}

	return violations
}
