package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github/BHAV0207/order-service/pkg/models"
	"io"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
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

/*
🧩 resp.Body — what it actually is
When you make an HTTP request like:
resp, err := http.Get(url)
The returned resp is of type *http.Response, and its Body field is an io.ReadCloser — meaning:
You can read from it (like a stream of bytes).
You must close it when you’re done.


⚙️ Why you must call defer resp.Body.Close()
Because every HTTP response keeps a network connection (TCP socket) open until you close the body.
If you don’t close it:
The connection stays open in the pool.
Eventually you’ll run out of file descriptors or sockets.
Future HTTP calls can hang or fail with too many open files or connection reset errors.
*/

func ValidateProduct(id primitive.ObjectID, price float64) bool {
	url := fmt.Sprintf("http://product-service:4000/get/%s", id.Hex())
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("Error contacting the productService", err)
		return false
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	} else {
		fmt.Println("product Service response body:", string(body))
	}
	fmt.Println(body)

	//cross checking the price logic is left

	fmt.Println("product Service status code:", resp.StatusCode)
	return resp.StatusCode == http.StatusOK
}

func ValidateInventory(id primitive.ObjectID) bool {
	url := fmt.Sprintf("http://inventory-service:6000/get/%s", id.Hex())
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("Error contacting inventory service", err)
		return false
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	} else {
		fmt.Println("product Service response body:", string(body))
	}

	fmt.Println(body)

	fmt.Println("product Service status code:", resp.StatusCode)
	return resp.StatusCode == http.StatusOK
}

func ValidateCartAndGetItems(id primitive.ObjectID) ([]map[string]interface{}, bool) {
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

func GetOrder(ctx context.Context, collection *mongo.Collection, id primitive.ObjectID) (models.Order, error) {
	var order models.Order

	filter := bson.M{"_id": id}

	err := collection.FindOne(ctx, filter).Decode(&order)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return order, errors.New("order not found")
		}
		return order, err
	}

	return order, nil
}

func GetAllOrderOfUser(ctx context.Context, collection *mongo.Collection, id primitive.ObjectID) ([]models.Order, error) {
	var orders []models.Order
	filter := bson.M{"userId": id}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &orders); err != nil {
		return nil, err
	}

	return orders, nil

}
