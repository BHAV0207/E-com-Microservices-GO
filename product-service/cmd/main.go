package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/BHAV0207/product-service/internal/handler"
	"github.com/BHAV0207/product-service/internal/repository"
	"github.com/gorilla/mux"
)

func main() {
	uri := "mongodb+srv://jainbhav0207_db_user:e94JtmF1QEGw0pxT@cluster0.uxo15jw.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"

	client := repository.ConnectDb(uri)
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	db := client.Database("ProductService")

	productHandler := &handler.ProductHandler{Collection: db.Collection("products")}

	router := mux.NewRouter()

	router.HandleFunc("/add", productHandler.CreateProduct).Methods("POST")
	router.HandleFunc("/delete/{id}", productHandler.DeleteProduct).Methods("DELETE")
	router.HandleFunc("/update/{id}", productHandler.UpdateProduct).Methods("PUT")
	router.HandleFunc("/get", productHandler.GetAllProducts).Methods("GET")
	router.HandleFunc("/get/{id}", productHandler.GetProductById).Methods("GET")
	router.HandleFunc("/get/user/{userId}", productHandler.GetProductsByUserId).Methods("GET")

	fmt.Println("Server listening on http://localhost:4000")
	log.Fatal(http.ListenAndServe(":4000", router))
}
