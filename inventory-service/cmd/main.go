package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BHAV0207/inventory-service/internal/event"
	"github.com/BHAV0207/inventory-service/internal/handler"
	"github.com/BHAV0207/inventory-service/internal/repository"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found — using system environment variables")
	}

	uri := os.Getenv("MONGO_INVENTORY_URI")
	if uri == "" {
		log.Fatal("MONGO_INVENTORY_URI is not set")
	}

	port := os.Getenv("PORT")

	client := repository.ConnectDb(uri)
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	db := client.Database("OrderService")
	inventoryCol := db.Collection("inventory")
	reservationCol := db.Collection("reservations")

	InventoryHandler := &handler.InventoryHandler{Collection: db.Collection("inventory")}
	router := mux.NewRouter()

	router.HandleFunc("/get/{id}", InventoryHandler.GetInventoryByProducId).Methods("GET")
	router.HandleFunc("/create", InventoryHandler.CreateInventory).Methods("POST")
	router.HandleFunc("/update/{id}", InventoryHandler.UpdateInventory).Methods("PUT")
	router.HandleFunc("/reserve", InventoryHandler.ReserveInventory).Methods("POST")

	go func() {
		kafkaBroker := os.Getenv("kafka:9092") // e.g. "localhost:9092"
		kafkaTopic := os.Getenv("payment-events") // e.g. "payment-events"
		kafkaGroup := "inventory-consumer-group"

		if kafkaBroker == "" || kafkaTopic == "" {
			log.Println("⚠️  KAFKA_BROKER or PAYMENT_TOPIC not set, skipping Kafka consumer")
			return
		}
		fmt.Printf("🧭 Starting Kafka consumer on topic '%s'...\n", kafkaTopic)
		event.ConsumePaymentEvents(kafkaBroker, kafkaTopic, kafkaGroup, inventoryCol, reservationCol)
	}()

	fmt.Printf("Server listening on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
