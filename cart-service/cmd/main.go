package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/BHAV0207/cart-service/internal/handler"
	"github.com/BHAV0207/cart-service/internal/repository"
	"github.com/gorilla/mux"
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

	cartHandelder := &handler.CartHandler{Collection: db.Collection("cart")}

	router := mux.NewRouter()

	router.HandleFunc("/addtocart", cartHandelder.AddToCart).Methods("POST")

	fmt.Println("Server listening on http://localhost:9000")
	log.Fatal(http.ListenAndServe(":9000", router))

}
