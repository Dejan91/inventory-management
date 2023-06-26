package gapi

import (
	"fmt"
	pb "github.com/Dejan91/inventory-management/proto/golang/user"
	"github.com/Dejan91/inventory-management/user/db"
	"github.com/Dejan91/inventory-management/user/store"
	"github.com/Dejan91/inventory-management/user/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type Server struct {
	pb.UnimplementedUserServer
	config *util.Config
	db     db.DB
	store  store.Store
}

func NewServer(config *util.Config, db db.DB, store store.Store) *Server {
	return &Server{
		db:     db,
		config: config,
		store:  store,
	}
}

func (s *Server) Run() error {
	grpcServer := grpc.NewServer()
	pb.RegisterUserServer(grpcServer, s)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", s.config.GrpcServerAddress)
	if err != nil {
		return err
	}

	fmt.Printf("starting gRPC server at %s\n", s.config.GrpcServerAddress)
	err = grpcServer.Serve(listener)
	if err != nil {
		return err
	}

	return nil
}
