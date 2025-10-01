package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/BHAV0207/user-service/internal/service"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	users, err := service.GetAll(ctx, h.Collection)
	if err != nil {
		http.Error(w, "Failed to fetch All users: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json") // Tell client: "I’m sending JSON"
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	//  mux.Vars(r) extracts the URL parameters from the incoming request (r).
	// It returns a map[string]string.
	vars := mux.Vars(r)
	idHex := vars["id"]

	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := service.GetById(ctx, h.Collection, id)
	if err != nil {
		http.Error(w, "Failed to get a user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json") // Tell client: "I’m sending JSON"
	json.NewEncoder(w).Encode(user)

	// json.NewEncoder(w).Encode(products)
	// json.NewEncoder(w) → creates a JSON encoder that writes directly to the response w.
	// .Encode(products) → takes your products (a slice or struct) and converts it into JSON format, then sends it into the response body.

}
func (h *UserHandler) UpdateUserDetails(w http.ResponseWriter, r *http.Request) {
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	modifiedCnt, err := service.Update(ctx, h.Collection, id, updateFields)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Updated %d user(s)", modifiedCnt)

}
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idHex := vars["id"]

	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	delCnt, err := service.Delete(ctx, h.Collection, id)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Deleted %d user(s)", delCnt)
}
