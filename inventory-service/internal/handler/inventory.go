package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/BHAV0207/inventory-service/internal/service"
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
