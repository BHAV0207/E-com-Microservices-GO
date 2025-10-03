package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/BHAV0207/product-service/internal/service"
	"github.com/BHAV0207/product-service/pkg/models"
	"github.com/BHAV0207/product-service/pkg/types"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	go func(productId primitive.ObjectID){
		if err := service.CreateInventoryForProduct(productId) ; err!= nil {
			fmt.Printf("⚠️ Failed to create inventory for product %s: %v\n", productId.Hex(), err)
		}
	}(id)


	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Inserted product with ID: %v", id)

}

func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	products, err := service.GetAll(ctx, h.Collection)
	if err != nil {
		http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json") // Tell client: "I’m sending JSON"
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) GetProductById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idHex := vars["id"]

	id, err := primitive.ObjectIDFromHex(idHex)
	fmt.Println(id)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var (
		product    models.Product
		inventory  types.InventoryResponse
		productErr error
		invErr     error
	)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		product, productErr = service.GetById(ctx, h.Collection, id)
	}()

	go func() {
		defer wg.Done()
		inventory, invErr = service.FetchInventoryByProductID(ctx, id)
	}()

	wg.Wait()

	// 🛑 Handle errors individually
	if productErr != nil {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}
	if invErr != nil {
		// You can choose to fail here OR continue without inventory
		http.Error(w, invErr.Error(), http.StatusBadGateway)
		return
	}
	response := struct {
		Product   models.Product   `json:"product"`
		Inventory types.InventoryResponse `json:"inventory"`
	}{
		Product:   product,
		Inventory: inventory,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *ProductHandler) GetProductsByUserId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idHex := vars["userId"]
	fmt.Println(idHex)
	id, err := primitive.ObjectIDFromHex(idHex)
	fmt.Println(id)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	products, err := service.GetByUserId(ctx, h.Collection, id)
	if err != nil {
		http.Error(w, "Failed to fetch products: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json") // Tell client: "I’m sending JSON"
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idHex := vars["id"]
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var updateFields map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updateFields); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	updatedFieldCount, err := service.UpdateProduct(context, h.Collection, id, updateFields)
	if err != nil {
		http.Error(w, "Failed to update product", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Updated %d product(s)", updatedFieldCount)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idHex := vars["id"]

	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		http.Error(w, "Failed to fetch id", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	delCnt, err := service.DeleteProduct(ctx, h.Collection, id)
	if err != nil {
		http.Error(w, "failed to delete the product", http.StatusInternalServerError)
	}

	fmt.Fprintf(w, "Deleted %d product(s)", delCnt)

}
