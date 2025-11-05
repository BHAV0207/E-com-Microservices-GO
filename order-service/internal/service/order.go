package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github/BHAV0207/order-service/pkg/models"
	"io"
	"math"
	"net/http"
	"time"

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
	fmt.Println(price, "from service file")

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

	var data map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println("❌ Error parsing product JSON:", err)
		return false
	}

	productData, ok := data["product"].(map[string]interface{})
	if !ok {
		fmt.Println("❌ product key missing or invalid")
		return false
	}

	priceVal, ok := productData["price"].(float64)
	if !ok {
		fmt.Println("❌ price value missing or invalid type")
		return false
	}

	const epsilon = 0.01
	if math.Abs(priceVal-price) > epsilon {
		fmt.Printf("❌ Price mismatch: expected %.2f, got %.2f\n", priceVal, price)
		return false
	}

	fmt.Println("product Service status code:", resp.StatusCode)
	return resp.StatusCode == http.StatusOK
}

func ValidateInventory(id primitive.ObjectID, quantity int64) bool {
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

	var data map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println("❌ Error parsing product JSON:", err)
		return false
	}

	stock, ok := data["inventory"].(float64)
	if !ok {
		fmt.Println("❌ stock value missing or invalid type")
		return false
	}

	if stock < float64(quantity) {
		fmt.Println("insuffecient stock")
		return false
	}
	fmt.Println(data)

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

// CallPaymentService sends payment initiation request from Order → Payment Service (HTTP)
func CallPaymentService(paymentReq map[string]any) error {
	url := "http://payment-service:3000/payments" // 👈 Adjust to your Payment Service URL

	jsonData, err := json.Marshal(paymentReq)
	if err != nil {
		return fmt.Errorf("failed to marshal payment request: %v", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("payment service call failed: %v", err)
	}
	defer resp.Body.Close()

	// If the Payment Service doesn’t return 201 or 200, treat as error
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return errors.New("payment service returned non-success: " + resp.Status)
	}

	fmt.Println("💳 Payment initiated successfully with Payment Service")
	return nil
}

// CallInventoryReserveAPI sends inventory reservation request from Order → Inventory Service (HTTP)
func CallInventoryReserveAPI(orderId string, items []map[string]any) (string, error) {
	url := "http://inventory-service:6000/api/inventory/reserve"

	body := map[string]any{
		"orderId": orderId,
		"items":   items,
	}

	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("inventory service not reachable: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("inventory reserve failed with status: %v", res.Status)
	}

	var resBody map[string]any
	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		return "", fmt.Errorf("invalid response from inventory")
	}

	reservationId, _ := resBody["reservationId"].(string)
	return reservationId, nil
}
