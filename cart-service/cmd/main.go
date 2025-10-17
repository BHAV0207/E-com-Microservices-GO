package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BHAV0207/cart-service/internal/events"
	"github.com/BHAV0207/cart-service/internal/handler"
	"github.com/BHAV0207/cart-service/internal/repository"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found — using system environment variables")
	}

	uri := os.Getenv("MONGO_CART_URI")
	if uri == "" {
		log.Fatal("MONGO_CART_URI is not set")
	}

	port := os.Getenv("PORT")

	client := repository.ConnectDb(uri)
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	db := client.Database("CartService")
	cartHandelder := &handler.CartHandler{Collection: db.Collection("cart")}

	// -------------------------------
	// Start Kafka Consumer in goroutine
	// -------------------------------
	go func() {
		consumer := events.NewConsumer(
			"kafka:9092", // e.g., "kafka:9092"
			"user-created",            // topic
			"cart-service-group",      // consumer group
			db.Collection("cart"),     // Mongo collection
		)
		consumer.Consume()
	}()

	router := mux.NewRouter()
	router.HandleFunc("/addtocart", cartHandelder.AddToCart).Methods("POST")
	router.HandleFunc("/user/{id}", cartHandelder.GetUsersCartById).Methods("GET")

	fmt.Println("Server listening on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
