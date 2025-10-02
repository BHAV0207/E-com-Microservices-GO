package main

import (
	"context"
	"fmt"
	"github/BHAV0207/order-service/internal/handler"
	"github/BHAV0207/order-service/internal/repository"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	uri := "mongodb+srv://jainbhav0207_db_user:WHMJ524qrJW27rDW@cluster0.wy3eykv.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"

	client := repository.ConnectDb(uri)
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	db := client.Database("OrderService")

	OrderHnadler := &handler.OrderHnadler{Collection: db.Collection("order")}
	router := mux.NewRouter()

	router.HandleFunc("/order", OrderHnadler.CreateOrder).Methods("POST")
	router.HandleFunc("/order/{id}", OrderHnadler.GetOrderByOrderId).Methods("GET")
	router.HandleFunc("/order/user/{id}", OrderHnadler.GetAllOrdersOfUserById).Methods("GET");

	fmt.Println("Server listening on http://localhost:7000")
	log.Fatal(http.ListenAndServe(":7000", router))

}
