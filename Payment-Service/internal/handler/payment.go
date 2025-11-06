package handler

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/BHAV0207/Product-service/internal/event"
	"github.com/BHAV0207/Product-service/internal/service"
	"github.com/BHAV0207/Product-service/pkg/models"
	"github.com/goccy/go-json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PaymentHandler struct {
	Collection *mongo.Collection
	Producer   *event.Producer
}

type OrderRequestBody struct {
	OrderID       string  `json:"orderId"`
	UserID        string  `json:"userId"`
	ReservationID string  `json:"reservationId"`
	Amount        float64 `json:"amount"`
	Method        string  `json:"method"`
}

func (h *PaymentHandler) PaymentCreation(w http.ResponseWriter, r *http.Request) {
	var reqBody OrderRequestBody

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	var userValid, orderValid bool

	wg.Add(2)

	// 🔹 Run user validation concurrently
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				userValid = false
			}
		}()
		userValid = service.ValidateUser(ctx, reqBody.UserID)
	}()

	// 🔹 Run order validation concurrently
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				orderValid = false
			}
		}()
		orderValid = service.ValidateOrder(ctx, reqBody.OrderID)
	}()

	// Wait for both to finish
	wg.Wait()

	payment := models.Paments{
		Id:            primitive.NewObjectID(),
		OrderID:       reqBody.OrderID,
		UserID:        reqBody.UserID,
		Amount:        reqBody.Amount,
		Currency:      "USD",
		Status:        "Failure",
		Method:        reqBody.Method,
		GatewayTxnID:  "TXN123456789",
		FailureReason: "lavde lag gaye",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if !userValid {
		h.Collection.InsertOne(ctx, payment)
		http.Error(w, "User validation failed", http.StatusBadRequest)
		return
	}

	if !orderValid {
		h.Collection.InsertOne(ctx, payment)
		http.Error(w, "Order validation failed", http.StatusBadRequest)
		return
	}

	// 🔹 Process payment after both validations succeed
	paymentSuccess := service.ProcessPayment(ctx, reqBody.OrderID, reqBody.Amount, reqBody.Method)
	if !paymentSuccess {
		h.Collection.InsertOne(ctx, payment)
		_ = h.Producer.Publish(event.PaymentEvent{
			OrderID:       reqBody.OrderID,
			UserID:        reqBody.UserID,
			ReservationID: reqBody.ReservationID,
			Amount:        reqBody.Amount,
			Method:        reqBody.Method,
			Status:        "failure",
			Timestamp:     time.Now(),
		})
		http.Error(w, "Payment processing failed", http.StatusInternalServerError)
		return
	}

	payment.Status = "success"
	_, err := h.Collection.InsertOne(ctx, payment)
	if err != nil {
		http.Error(w, "Failed to save payment to database", http.StatusInternalServerError)
		return
	}

	_ = h.Producer.Publish(event.PaymentEvent{
		OrderID:       reqBody.OrderID,
		UserID:        reqBody.UserID,
		ReservationID: reqBody.ReservationID,
		Amount:        reqBody.Amount,
		Method:        reqBody.Method,
		Status:        "success",
		Timestamp:     time.Now(),
	})

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Payment processed successfully"))
}
