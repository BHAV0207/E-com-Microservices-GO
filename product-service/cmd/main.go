package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BHAV0207/product-service/internal/handler"
	"github.com/BHAV0207/product-service/internal/repository"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	println(godotenv.Load())

	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  No .env file found — using system environment variables")
	}

	uri := os.Getenv("MONGO_PRODUCT_URI")
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

	db := client.Database("ProductService")

	productHandler := &handler.ProductHandler{Collection: db.Collection("products")}

	router := mux.NewRouter()

	router.HandleFunc("/add", productHandler.CreateProduct).Methods("POST")
	router.HandleFunc("/delete/{id}", productHandler.DeleteProduct).Methods("DELETE")
	router.HandleFunc("/update/{id}", productHandler.UpdateProduct).Methods("PUT")
	router.HandleFunc("/get", productHandler.GetAllProducts).Methods("GET")
	router.HandleFunc("/get/{id}", productHandler.GetProductById).Methods("GET")
	router.HandleFunc("/get/user/{userId}", productHandler.GetProductsByUserId).Methods("GET")

	fmt.Println("Server listening on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
