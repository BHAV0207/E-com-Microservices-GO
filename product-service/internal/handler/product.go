package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/BHAV0207/product-service/internal/service"
	"github.com/BHAV0207/product-service/pkg/models"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductHandler struct {
	Collection *mongo.Collection
}

var validate = validator.New()

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {

	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid rewuest body", http.StatusBadRequest)
	}

	// Validate struct fields
	// checks for the required fields
	if err := validate.Struct(product); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, err := service.InsertProduct(ctx, h.Collection, product)
	if err != nil {
		http.Error(w, "Failed to insert product", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Inserted product with ID: %v", id)

}
