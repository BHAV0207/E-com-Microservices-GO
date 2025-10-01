package main

import (
	"context"
	"github/BHAV0207/order-service/internal/repository"
	"log"
)

func main() {
	uri := "mongodb+srv://jainbhav0207_db_user:WHMJ524qrJW27rDW@cluster0.wy3eykv.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"

	client := repository.ConnectDb(uri)
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	db := client.Database("OrderService");

	

}
