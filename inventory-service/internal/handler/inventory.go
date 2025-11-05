package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/BHAV0207/inventory-service/internal/service"
	"github.com/BHAV0207/inventory-service/pkg/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type InventoryHandler struct {
	Collection *mongo.Collection
}

func (h *InventoryHandler) GetInventoryByProducId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idHex := vars["id"]
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		http.Error(w, "id invalid", http.StatusBadRequest)
		return
	}

	filter := bson.M{"productId": id}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	inventory, err := service.Get(ctx, h.Collection, filter)
	if err != nil {
		http.Error(w, "inventory not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(inventory)
}

type InventoryRequest struct {
	ID        primitive.ObjectID `json:"_id"`
	ProductId primitive.ObjectID `json:"productId"`
	Inventory int64              `json:"inventory"`
}

func (h *InventoryHandler) CreateInventory(w http.ResponseWriter, r *http.Request) {
	var req InventoryRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body format provided", http.StatusBadRequest)
		return
	}

	if req.ProductId.IsZero() {
		http.Error(w, "ProductId is required", http.StatusBadRequest)
		return
	}

	if req.Inventory < 0 {
		http.Error(w, "Inventory cannot be negative", http.StatusBadRequest)
		return
	}

	inventory := models.Inventory{
		ID:        primitive.NewObjectID(),
		ProductId: req.ProductId,
		Inventory: req.Inventory,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := service.Create(ctx, h.Collection, inventory)
	if err != nil {
		http.Error(w, "Failed to create inventory: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "Inventory created successfully",
		"inventory":  inventory,
		"insertedId": result.InsertedID,
	})

}

func (h *InventoryHandler) UpdateInventory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idHex := vars["id"]
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var updateFields map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updateFields); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	updateFields["updatedAt"] = time.Now()

	// 4️⃣ Context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 5️⃣ Call service layer
	updatedInventory, err := service.Update(ctx, h.Collection, id, updateFields)
	if err != nil {
		http.Error(w, "failed to update inventory: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 6️⃣ Return updated document as response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedInventory)
}

func (h *InventoryHandler) DeleteInventory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idHex := vars["id"]
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		http.Error(w, "id not valid ", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	delCnt, err := service.DeleteProduct(ctx, h.Collection, id)
	if err != nil {
		http.Error(w, "failed to delete the product", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Deleted %d product(s)", delCnt)

}

func (h *InventoryHandler) ReserveInventory(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OrderId string                   `json:"orderId"`
		Items   []models.ReservationItem `json:"items"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.OrderId == "" || len(req.Items) == 0 {
		http.Error(w, "Missing orderId or items", http.StatusBadRequest)
		return
	}

	context, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for _, item := range req.Items {
		prodId, _ := primitive.ObjectIDFromHex(item.ProductId)
		var inventory models.Inventory

		err := h.Collection.FindOne(context, bson.M{"productId": prodId}).Decode(&inventory)
		if err != nil {
			http.Error(w, "Product not found: "+item.ProductId, http.StatusNotFound)
			return
		}

		if inventory.Inventory < item.Quantity {
			http.Error(w, "Insufficient inventory for product: "+item.ProductId, http.StatusConflict)
			return
		}
	}

	reservation := models.Reservation{
		ID:            primitive.NewObjectID(),
		ReservationID: primitive.NewObjectID().Hex(),
		OrderID:       req.OrderId,
		Status:        "PENDING",
		Items:         req.Items,
		ExpiresAt:     time.Now().Add(15 * time.Minute),
		CreatedAt:     time.Now(),
	}

	resCol := h.Collection.Database().Collection("reservations")
	_, err := resCol.InsertOne(context, reservation)
	if err != nil {
		http.Error(w, "Failed to create reservation: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"success":       true,
		"reservationId": reservation.ReservationID,
		"expiresAt":     reservation.ExpiresAt,
		"message":       "Reservation created successfully",
	})

}
