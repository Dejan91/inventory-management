package main

import (
	"context"
	"fmt"
	"github.com/Dejan91/inventory-management/user/db"
	"github.com/Dejan91/inventory-management/user/model"
	"github.com/Dejan91/inventory-management/user/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	fmt.Println("user service")
	mongoDB, err := db.NewMongoDB(config.DBUser, config.DBPassword, config.DBHost, config.DBName)
	if err != nil {
		panic(err)
	}

	user, err := mongoDB.Create(context.Background(), &model.User{
		Username: "dejan",
		Email:    "test@test.com",
	})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(user)
}
