package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/BHAV0207/user-service/internal/events"
	"github.com/BHAV0207/user-service/internal/service"
	"github.com/BHAV0207/user-service/pkg/models"
	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt"
)

var validate = validator.New()

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if err := validate.Struct(user); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	hashedPass, err := service.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "error in hasing password", http.StatusBadRequest)
	}

	user.Password = hashedPass
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// ✅ FIX: Capture all three return values
	_, exists, err := service.GetUserByEmail(ctx, h.Collection, user.Email)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if exists {
		http.Error(w, "User already exists, please login", http.StatusConflict)
		return
	}

	res, err := h.Collection.InsertOne(ctx, user)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	producer := events.NewProducer("kafka:9092", "user-created")
	event := map[string]interface{}{
		"userId": res.InsertedID,
		"email":  user.Email,
	}
	if err := producer.Publish(event); err != nil {
		log.Printf("⚠️ Failed to publish user-created event: %v", err)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User registered successfully",
	})
}

// LOGIN -->

var jwtKey = []byte("your_secret_key")

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `josn:"token"`
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, exists, err := service.GetUserByEmail(ctx, h.Collection, req.Email)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !exists {
		http.Error(w, "User does not  exists, please register", http.StatusUnauthorized)
		return
	}

	if !service.ComparePasswordHash(req.Password, user.Password) {
		http.Error(w, "Invalid password for hashing ", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{Token: tokenString})
}
