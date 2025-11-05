package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github/BHAV0207/order-service/internal/event"
	"github/BHAV0207/order-service/internal/handler"
	"github/BHAV0207/order-service/internal/repository"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found — using system environment variables")
	}

	uri := os.Getenv("MONGO_ORDER_URI")
	if uri == "" {
		log.Fatal("MONGO_ORDER_URI is not set")
	}

	port := os.Getenv("PORT")
	// if port == "" {
	// 	port = "7000"
	// }

	client := repository.ConnectDb(uri)
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	db := client.Database("OrderService")
	OrderHnadler := &handler.OrderHnadler{Collection: db.Collection("order")}

	go func() {
		consumer := event.NewConsumer(
			"kafka:9092",
			"payment-events",
			"order-service-payment-group",
			db.Collection("order"),
		)
		consumer.ListenForPayments()
	}()

	router := mux.NewRouter()

	router.HandleFunc("/order", OrderHnadler.CreateOrder).Methods("POST")
	router.HandleFunc("/order/{id}", OrderHnadler.GetOrderByOrderId).Methods("GET")
	router.HandleFunc("/order/user/{id}", OrderHnadler.GetAllOrdersOfUserById).Methods("GET")

	fmt.Println("Server listening on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
