package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/BHAV0207/user-service/internal/service"
	"github.com/BHAV0207/user-service/pkg/models"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	Collection *mongo.Collection
}

var validate = validator.New()

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "bad reuquest", http.StatusInternalServerError)
	}

	if err := validate.Struct(user); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exists, err := service.GetUserByEmail(ctx, h.Collection, user.Email)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if exists {
		http.Error(w, "User already exists, please login", http.StatusConflict) // 409 Conflict
		return
	}

	_, err = h.Collection.InsertOne(ctx, user)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User registered Successfully",
	})
}
