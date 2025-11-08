package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/BHAV0207/payment-service/internal/event"
	"github.com/BHAV0207/payment-service/internal/repository"
	"github.com/joho/godotenv"
)

func main() {
	// 1️⃣ Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found — using system environment variables")
	}

	mongoURI := os.Getenv("MONGO_NOTIFICATION_URI")
	port := os.Getenv("PORT")

	// 2️⃣ Validate configuration
	if mongoURI == "" {
		log.Fatal("❌ MONGO_NOTIFICATION_URI not set in environment")
	}
	if port == "" {
		port = "2000"
	}

	// 3️⃣ Connect to MongoDB
	client := repository.ConnectDb(mongoURI)
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	db := client.Database("NotificationService")
	notifCol := db.Collection("notifications")

	// 4️⃣ Create and start Kafka consumers for multiple topics
	// orderConsumer := event.NewConsumer(
	// 	"kafka:9092",
	// 	orderTopic,
	// 	"notif-order-group",
	// 	"OrderConsumer",
	// 	notifCol,
	// )
	paymentConsumer := event.NewConsumer(
		"kafka:9092",
		"payment-events",
		"notif-payment-group",
		"payment-service",
		notifCol,
	)

	userConsumer := event.NewConsumer(
		"kafka:9092",
		"user-created",
		"notif-user-group",
		"user-service",
		notifCol,
	)

	userDeletedConsumer := event.NewConsumer(
		"kafka:9092",
		"user-deleted",
		"notif-userdel-group",
		"user-service",
		notifCol,
	)
	
	// Run both consumers concurrently
	// go orderConsumer.StartConsuming()
	go userDeletedConsumer.StartConsuming()
	go userConsumer.StartConsuming()
	go paymentConsumer.StartConsuming()

	// 5️⃣ Start HTTP server (optional for future APIs)
	server := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	fmt.Printf("🔔 Notification Service running on http://localhost:%s\n", port)
	log.Fatal(server.ListenAndServe())
}
