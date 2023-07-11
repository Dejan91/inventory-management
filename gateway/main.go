package main

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"github.com/Dejan91/inventory-management/gateway/client"
	"github.com/Dejan91/inventory-management/gateway/server"
	"github.com/Dejan91/inventory-management/gateway/service"
	"github.com/Dejan91/inventory-management/gateway/util"
	pb "github.com/Dejan91/inventory-management/proto/gateway/user/api/v1"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

const (
	listenAddress = "0.0.0.0:9090"
)

func main() {
	ctx := context.Background()

	log.Printf("Gateway server starting on %s", listenAddress)

	config, err := util.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	opt := option.WithCredentialsFile("credentials/inventory-management-1f296-firebase-adminsdk-cre40-d6530cc61f.json")
	conf := &firebase.Config{ProjectID: config.FirebaseProjectID}
	firebaseApp, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalf("error initializing Firebase app: %v\n", err)
	}

	authClient, err := firebaseApp.Auth(context.Background())
	if err != nil {
		log.Fatalf("failed to create Firebase auth client: %v", err)
	}

	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	authInterceptor := service.NewAuthInterceptor(authClient)
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.Unary()),
	)
	reflection.Register(grpcServer)

	userClient, err := client.NewUserClient(ctx)
	if err != nil {
		panic(err)
	}

	pb.RegisterUserServer(grpcServer, server.NewUserServer(userClient))

	if err = grpcServer.Serve(listener); err != nil {
		log.Fatalf("gRPC server failed: %v", err)
	}
}
