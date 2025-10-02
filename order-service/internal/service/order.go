package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ValidateUser(id primitive.ObjectID) bool {
	fmt.Println(id)
	url := fmt.Sprintf("http://user-service:8080/users/%s", id.Hex())
	fmt.Println("Calling User service with URL:", url)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error contacting User-service:", err)
		return false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	} else {
		fmt.Println("User Service response body:", string(body))
	}

	fmt.Println("User Service status code:", resp.StatusCode)
	return resp.StatusCode == http.StatusOK
}

func ValidateCartAndGetItems(id primitive.ObjectID) (interface{}, bool) {
	url := fmt.Sprintf("http://cart-service:9000/user/%s", id.Hex())
	fmt.Println("Calling Cart service with URL:", url)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error contacting Cart-service:", err)
		return nil, false
	}
	defer resp.Body.Close()

	// If cart not found
	if resp.StatusCode != http.StatusOK {
		fmt.Println("❌ Cart service returned non-200 status:", resp.StatusCode)
		return nil, false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	} else {
		fmt.Println("CART Service response body:", string(body))
	}

	var items []map[string]interface{}
	if err := json.Unmarshal(body, &items); err != nil {
		fmt.Println("❌ Error parsing cart JSON:", err)
		return nil, false
	}

	// .([]interface{}) is a type assertion in Go. It converts a generic interface{} (from JSON) into a slice []interface{}, allowing you to loop through the array.
	if len(items) == 0 {
		fmt.Println("⚠️ No items found in cart")
		return nil, false
	}

	return items, true

}
