package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/BHAV0207/Product-service/internal/event"
	"github.com/BHAV0207/Product-service/internal/handler"
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

	db := client.Database("PaymentService")
	collection := db.Collection("payments")

	producer := event.NewProducer("kafka:9092", "payment-events")

	paymentHandler := &handler.PaymentHandler{
		Collection: collection,
		Producer:   producer,
	}
	http.HandleFunc("/payments", paymentHandler.PaymentCreation)

	log.Println("🚀 Payment Service running on :3000")
	http.ListenAndServe(":3000", nil)
}
