package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/BHAV0207/cart-service/pkg"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

func ValidateUser(id string) bool {
	url := fmt.Sprintf("http://user-service:8080/users/%s", id)
	fmt.Println("Calling User Service with URL:", url)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error contacting user-service:", err)
		return false
	}
	defer resp.Body.Close()

	// Read the response body for debugging
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	} else {
		fmt.Println("User Service response body:", string(body))
	}

	fmt.Println("User Service status code:", resp.StatusCode)
	return resp.StatusCode == http.StatusOK
}

func ValidateProduct(id string, quantity int64) bool {
	url := fmt.Sprintf("http://product-service:4000/get/%s", id)
	fmt.Println("Calling Product Service with URL:", url)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error contacting product-service:", err)
		return false
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("❌ Error reading product service response body:", err)
		return false
	}
	fmt.Println("✅ Product Service response body:", string(body))

	var data map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println("❌ Error parsing JSON:", err)
		return false
	}

	availableQty := int64(data["inventory"].(float64))

	if availableQty < quantity {
		fmt.Println("❌ Not enough stock available")
		return false
	}

	fmt.Println("✅ Sufficient stock available")

	fmt.Println("Product Service status code:", resp.StatusCode)
	return resp.StatusCode == http.StatusOK
}

// AddItemToCart adds a product to a user's cart or updates quantity if it already exists
func AddItemToCart(ctx context.Context, collection *mongo.Collection, userId, productId primitive.ObjectID, quantity int64) error {
	var cart pkg.Cart

	// 1. Try to find existing cart
	err := collection.FindOne(ctx, bson.M{"userId": userId}).Decode(&cart)
	if err == mongo.ErrNoDocuments {
		// 2. No cart exists, create a new one
		newCart := pkg.Cart{
			UserId: userId,
			Items: []pkg.CartItem{
				{ProductId: productId, Quantity: quantity},
			},
		}
		_, err := collection.InsertOne(ctx, newCart)
		return err
	} else if err != nil {
		return err
	}

	// 3. Cart exists: check if product already in items
	found := false
	for i, item := range cart.Items {
		if item.ProductId == productId {
			cart.Items[i].Quantity += quantity // update quantity
			found = true
			break
		}
	}

	// 4. If product not found, append new item
	if !found {
		cart.Items = append(cart.Items, pkg.CartItem{ProductId: productId, Quantity: quantity})
	}

	// 5. Update the cart in MongoDB
	_, err = collection.UpdateOne(
		ctx,
		bson.M{"userId": userId}, // use correct field name
		bson.M{"$set": bson.M{"items": cart.Items}},
	)
	return err
}

func GetUserCart(ctx context.Context, collection *mongo.Collection, userId primitive.ObjectID) (interface{}, error) {
	var cart pkg.Cart
	err := collection.FindOne(ctx, bson.M{"userId": userId}).Decode(&cart)
	if err != nil {
		fmt.Println("❌ Error fetching cart from DB:", err)
		return nil, err
	}

	var expandedCart []map[string]interface{}

	for _, item := range cart.Items {
		productId := item.ProductId.Hex()
		productUrl := fmt.Sprintf("http://product-service:4000/get/%s", productId)

		resp, err := http.Get(productUrl)
		if err != nil {
			fmt.Println("❌ Error contacting product service:", err)
			continue // skip this product
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("❌ Error reading product response:", err)
			continue
		}

		var data map[string]interface{}
		if err := json.Unmarshal(body, &data); err != nil {
			fmt.Println("❌ Error parsing product JSON:", err)
			continue
		}

		data["productId"] = productId
		data["quantity"] = item.Quantity

		expandedCart = append(expandedCart, data)
	}

	return expandedCart, nil
}

func CreateCart(ctx context.Context, collection *mongo.Collection, cart pkg.Cart) error {
	_, err := collection.InsertOne(ctx, cart)
	return err
}

func DeleteCart(ctx context.Context, collection *mongo.Collection, userID primitive.ObjectID) (int, error) {
	result, err := collection.DeleteOne(ctx, bson.M{"userId": userID})
	if err != nil {
		return 0, err
	}

	return int(result.DeletedCount), nil
}
