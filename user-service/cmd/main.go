package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/BHAV0207/user-service/internal/handler"
	"github.com/BHAV0207/user-service/internal/repository"
	"github.com/gorilla/mux"
)

func main() {
	uri := "mongodb+srv://jainbhav0207_db_user:XB9P4Jgp0fzqBCOS@cluster0.oa5vscu.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"
	client := repository.ConnectDb(uri)
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	db := client.Database("UserService")

	UserHandler := &handler.UserHandler{Collection: db.Collection("users")}

	router := mux.NewRouter()

	router.HandleFunc("/register", UserHandler.Register).Methods("POST")
	router.HandleFunc("/login", UserHandler.Login).Methods("POST")

	fmt.Println("Server listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))

}
