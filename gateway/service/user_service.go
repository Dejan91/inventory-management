package service

import (
	pb "github.com/Dejan91/inventory-management/proto/gateway/user/api/v1"
	userSvcV1 "github.com/Dejan91/inventory-management/proto/user/api/v1"
)

type userService struct {
	usersClient userSvcV1.UserClient
}
