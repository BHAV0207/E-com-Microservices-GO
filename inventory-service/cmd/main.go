package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

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
	InventoryHandler := &handler.InventoryHandler{Collection: db.Collection("inventory")}

	router := mux.NewRouter()
	router.HandleFunc("/get/{id}", InventoryHandler.GetInventoryByProducId).Methods("GET")
	router.HandleFunc("/create", InventoryHandler.CreateInventory).Methods("POST")
	router.HandleFunc("/update/{id}", InventoryHandler.UpdateInventory).Methods("PUT")

	fmt.Println("Server listening on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
