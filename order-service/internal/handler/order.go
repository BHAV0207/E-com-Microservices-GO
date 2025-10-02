package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github/BHAV0207/order-service/internal/service"
	"github/BHAV0207/order-service/pkg/models"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderHnadler struct {
	Collection *mongo.Collection
}

type CreateOrderRequest struct {
	UserId      string `json:"userId"`
	CartId      string `json:"cartId"`
	Address     string `json:"address"`
	PaymentInfo string `json:"paymentInfo"`
}

func (h *OrderHnadler) CreateOrder(w http.ResponseWriter, r *http.Request) {

	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "failed to process and get the user data", http.StatusBadRequest)
		return
	}
	fmt.Println("📦 Received Order Data:", req)

	userId, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		http.Error(w, "invalid userId", http.StatusBadRequest)
		return
	}

	if !service.ValidateUser(userId) {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// 2️⃣ Validate cart and get items
	cartItemsSlice, ok := service.ValidateCartAndGetItems(userId)
	if !ok {
		http.Error(w, "Cart is empty or not found", http.StatusBadRequest)
		return
	}

	var orderItems []models.OrderItems
	var total float64 = 0

	for _, it := range cartItemsSlice {
		// Convert productId
		productIdStr, _ := it["productId"].(string)
		productId, _ := primitive.ObjectIDFromHex(productIdStr)

		// Convert quantity
		quantityFloat, _ := it["quantity"].(float64)
		quantity := int64(quantityFloat)

		// Convert price
		priceFloat, _ := it["price"].(float64)

		orderItems = append(orderItems, models.OrderItems{
			ProductId: productId,
			Quantity:  quantity,
			Price:     priceFloat,
		})

		total += priceFloat * float64(quantity)
	}

	if len(orderItems) == 0 {
		http.Error(w, "No valid items in cart", http.StatusBadRequest)
		return
	}

	// 3️⃣ Get address from request
	address := req.Address
	if address == "" {
		http.Error(w, "Address is required", http.StatusBadRequest)
		return
	}

	// 4️⃣ Create order object
	order := models.Order{
		Id:        primitive.NewObjectID(),
		UserId:    userId,
		Items:     orderItems,
		Address:   address,
		Status:    "pending",
		CreatedAt: time.Now(),
		Total:     total,
	}

	// 5️⃣ Insert into MongoDB
	_, err = h.Collection.InsertOne(r.Context(), order)
	if err != nil {
		fmt.Println("❌ Error inserting order:", err)
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	// 6️⃣ Respond with created order
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHnadler) GetOrderByOrderId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idHex := vars["id"]

	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		http.Error(w, "invalid order ID format", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	order, err := service.GetOrder(ctx, h.Collection, id)
	if err != nil {
		http.Error(w, "couldn't fetch the order: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(order)
}
