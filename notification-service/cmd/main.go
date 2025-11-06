// package main

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"os"
// 	"time"
// 	VF4bHCwD8XXBJltb

// 	mongodb+srv://jainbhav0207_db_user:VF4bHCwD8XXBJltb@cluster0.ymwavs0.mongodb.net/?appName=Cluster0
// 	"github.com/BHAV0207/notification-service/internal/handler"
// 	"github.com/BHAV0207/notification-service/internal/repository"
// 	"github.com/BHAV0207/notification-service/kafka"
// 	"github.com/gorilla/mux"
// 	"github.com/joho/godotenv"
// )

// func main() {
// 	_ = godotenv.Load()
// 	mongoURI := os.Getenv("MONGO_URI")

// 	kafkaBroker := os.Getenv("KAFKA_BROKER")
// 	orderTopic := os.Getenv("ORDER_TOPIC")
// 	paymentTopic := os.Getenv("PAYMENT_TOPIC")
// 	port := os.Getenv("PORT")

// 	if port == "" {
// 		port = "8085"
// 	}

// 	client := repository.ConnectDb(mongoURI)
// 	db := client.Database("NotificationService")
// 	notifCol := db.Collection("notifications")

// 	// API routes
// 	router := mux.NewRouter()
// 	handler := &handler.NotificationHandler{Collection: notifCol}
// 	router.HandleFunc("/notifications", handler.GetUserNotifications).Methods("GET")

// 	// Start consumers for order & payment topics
// 	go kafka.StartKafkaConsumer(kafkaBroker, orderTopic, "notif-order-group", notifCol)
// 	go kafka.StartKafkaConsumer(kafkaBroker, paymentTopic, "notif-payment-group", notifCol)

// 	// Start HTTP server
// 	server := &http.Server{
// 		Addr:         ":" + port,
// 		Handler:      router,
// 		ReadTimeout:  15 * time.Second,
// 		WriteTimeout: 15 * time.Second,
// 		IdleTimeout:  60 * time.Second,
// 	}

// 	fmt.Printf("🔔 Notification Service running on http://localhost:%s\n", port)
// 	log.Fatal(server.ListenAndServe())
// }
