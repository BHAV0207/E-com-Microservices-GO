package main

import (
	"context"
	"log"
	"os"

	"github.com/BHAV0207/Product-service/internal/repository"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  No .env file found — using system environment variables")
	}

	uri := os.Getenv("MONGO_PAYMENT_URI")
	if uri == "" {
		log.Fatal("MONGO_PRODUCT_URI is not set")
	}
	println(uri)

	port := os.Getenv("PORT")
	println(port)

	client := repository.ConnectDb(uri)
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	// db := client.Database("PaymentService")
	}
