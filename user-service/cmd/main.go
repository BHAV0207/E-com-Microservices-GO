package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BHAV0207/user-service/internal/events"
	"github.com/BHAV0207/user-service/internal/handler"
	"github.com/BHAV0207/user-service/internal/repository"
	workerpool "github.com/BHAV0207/user-service/internal/workerPool"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var userCreatedPool *workerpool.WorkerPool
var userDeletedPool *workerpool.WorkerPool

func main() {
	/* Your worker pool setup issue

	   In your DeleteUser handler:

	   producer := events.NewProducer("kafka:9092", "user-deleted")
	   kafkaPool = workerpool.NewWorkerPool(10, producer)
	   log.Println("🚀 Kafka worker pool started with 10 workers")
	   event := map[string]any{"userId": id}
	   go kafkaPool.Submit(event)


	   Two issues here:
	   You are creating a new producer + new worker pool every time a user is deleted.
	   If 1 million deletes happen, you spawn 1 million worker pools, which defeats the purpose of controlling concurrency.
	   Correct approach: start the worker pool once at service startup, and reuse it for all events.
	   go kafkaPool.Submit(event) is unnecessary.
	   Submit only writes to the channel. The workers already read from it in goroutines.
	   Wrapping it in go is redundant.
	*/
  userCreatedProducer := events.NewProducer("kafka:9092", "user-created")
    userCreatedPool = workerpool.NewWorkerPool(10, userCreatedProducer)

    userDeletedProducer := events.NewProducer("kafka:9092", "user-deleted")
    userDeletedPool = workerpool.NewWorkerPool(10, userDeletedProducer)

 
	// ✅ Load environment variables
	err := godotenv.Load() // Only for local dev, safe to ignore in prod
	if err != nil {
		log.Println("⚠️  No .env file found — using system environment variables")
	}

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
	userHandler := &handler.UserHandler{
		Collection:       db.Collection("users"),
		UserCreatedPool:  userCreatedPool,
		UserDeletedPool:  userDeletedPool,
}


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
