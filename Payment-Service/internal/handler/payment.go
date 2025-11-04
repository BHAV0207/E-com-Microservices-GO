package handler

import (
	"net/http"

	"github.com/goccy/go-json"
	"go.mongodb.org/mongo-driver/mongo"
)

type PaymentHandler struct {
	Collection *mongo.Collection
}

type OrderRequestBody struct {
	orderId string
	userId  string
	amount  float64
	method  string
}

func (h *PaymentHandler) PaymentCreation(w http.ResponseWriter, r *http.Request) {
	var reqBody OrderRequestBody
	//  request body say data lena
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// validate user and order 
	
	//

	//  payment create hoga by payment model sey
	//

}
