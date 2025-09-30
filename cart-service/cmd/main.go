package main

import (
	"context"
	"log"

	"github.com/BHAV0207/cart-service/internal/repository"
)

func main() {
	uri := "mongodb+srv://jainbhav0207_db_user:t6rZLlfyCzTqPlex@cluster0.cgb8ri1.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"

	client := repository.ConnectDb(uri)
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	db := client.Database("CartService")

}
