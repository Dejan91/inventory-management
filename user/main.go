package main

import (
	"context"
	"github.com/Dejan91/inventory-management/user/db"
	"github.com/Dejan91/inventory-management/user/gapi"
	"github.com/Dejan91/inventory-management/user/store"
	"github.com/Dejan91/inventory-management/user/util"
	"google.golang.org/api/option"
	"log"

	firebase "firebase.google.com/go/v4"
)

func main() {
	ctx := context.Background()

	config, err := util.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	mongoDB, err := db.NewMongoDB(config.DBUser, config.DBPassword, config.DBHost, config.DBName)
	if err != nil {
		panic(err)
	}
	defer func(mongoDB db.DB, ctx context.Context) {
		err = mongoDB.Close(ctx)
		if err != nil {
			log.Fatal("error closing connection with mongodb")
		}
	}(mongoDB, ctx)

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
	firebaseStore := store.NewFirebaseStore(authClient)

	server := gapi.NewServer(config, mongoDB, firebaseStore)
	err = server.Run()
	if err != nil {
		log.Fatal("error running gRPC server", err)
	}
}
