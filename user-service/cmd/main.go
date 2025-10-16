package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BHAV0207/user-service/internal/handler"
	"github.com/BHAV0207/user-service/internal/repository"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	println(godotenv.Load())

	// ✅ Load environment variables
	err := godotenv.Load() // Only for local dev, safe to ignore in prod
	if err != nil {
		log.Println("⚠️  No .env file found — using system environment variables")
	}

	println(godotenv.Load())

	// ✅ Get env vars
	mongoURI := os.Getenv("MONGO_USER_URI")
	if mongoURI == "" {
		log.Fatal("❌ MONGO_URI not set in environment")
	}

	port := os.Getenv("PORT")
	println(port)
	println(mongoURI)

	// ✅ MongoDB connection
	client := repository.ConnectDb(mongoURI)
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
	router.HandleFunc("/users", userHandler.GetAllUsers).Methods("GET")            // ✅ Get all users
	router.HandleFunc("/users/{id}", userHandler.GetUserById).Methods("GET")       // ✅ Get user by ID
	router.HandleFunc("/users/{id}", userHandler.UpdateUserDetails).Methods("PUT") // ✅ Update user
	router.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")     // ✅ Delete user

	fmt.Println("🚀 Server running on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
