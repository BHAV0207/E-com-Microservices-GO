package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/BHAV0207/cart-service/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CartHandler struct {
	Collection *mongo.Collection
}

func (h *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Body)
	var req struct {
		UserID    string `json:"userId"`
		ProductID string `json:"productId"`
		Quantity  int64  `json:"quantity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userId, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		http.Error(w, "Invalid userId", http.StatusBadRequest)
		return
	}

	productId, err := primitive.ObjectIDFromHex(req.ProductID)
	if err != nil {
		http.Error(w, "Invalid productId", http.StatusBadRequest)
		return
	}
	if !service.ValidateUser(req.UserID) {
		fmt.Println("Validation failed for UserID:", req.UserID)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if !service.ValidateProduct(req.ProductID , req.Quantity) {
		fmt.Println("Validation failed for ProductID:", req.ProductID)
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = service.AddItemToCart(ctx, h.Collection, userId, productId, req.Quantity)
	if err != nil {
		http.Error(w, "Failed to add to cart", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Item added to cart"}`))

}
