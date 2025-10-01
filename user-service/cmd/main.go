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
	// ✅ MongoDB connection
	uri := "mongodb+srv://jainbhav0207_db_user:XB9P4Jgp0fzqBCOS@cluster0.oa5vscu.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"
	client := repository.ConnectDb(uri)
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	db := client.Database("UserService")
	userHandler := &handler.UserHandler{Collection: db.Collection("users")}

	// ✅ Router setup
	router := mux.NewRouter()

	// ------------------ AUTH ------------------
	router.HandleFunc("/register", userHandler.Register).Methods("POST")
	router.HandleFunc("/login", userHandler.Login).Methods("POST")

	// ------------------ CRUD USERS ------------------
	router.HandleFunc("/users", userHandler.GetAllUsers).Methods("GET")             // ✅ Get all users
	router.HandleFunc("/users/{id}", userHandler.GetUserById).Methods("GET")        // ✅ Get user by ID
	router.HandleFunc("/users/{id}", userHandler.UpdateUserDetails).Methods("PUT")  // ✅ Update user
	router.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")      // ✅ Delete user

	fmt.Println("🚀 Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
