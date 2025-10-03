package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/BHAV0207/product-service/pkg/models"
	"github.com/BHAV0207/product-service/pkg/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAll(ctx context.Context, collection *mongo.Collection) ([]models.Product, error) {
	var products []models.Product

	cursor, err := collection.Find(ctx, bson.M{}) // ✅ FIXED
	if err != nil {
		log.Printf("❌ MongoDB Find error: %v\n", err)
		return nil, err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var product models.Product
		if err := cursor.Decode(&product); err != nil {
			log.Printf("❌ Decode error: %v\n", err)
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func GetById(ctx context.Context, collection *mongo.Collection, id primitive.ObjectID) (models.Product, error) {
	var product models.Product

	// ✅ Use "_id" (not "id")
	filter := bson.M{"_id": id}

	// ✅ Use FindOne instead of Find (because we want only one document)
	err := collection.FindOne(ctx, filter).Decode(&product)
	if err != nil {
		return product, err // return zero-value product + error
	}

	return product, nil
}
func GetByUserId(ctx context.Context, collection *mongo.Collection, id primitive.ObjectID) ([]models.Product, error) {
	var products []models.Product

	cursor, err := collection.Find(ctx, bson.D{{Key: "userId", Value: id}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var product models.Product
		if err := cursor.Decode(&product); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func InsertProduct(ctx context.Context, collection *mongo.Collection, product models.Product) (primitive.ObjectID, error) {
	result, err := collection.InsertOne(ctx, product)
	if err != nil {
		return primitive.NilObjectID, err
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, fmt.Errorf("failed to assert InsertedID to ObjectID")
	}

	return id, nil
}

func UpdateProduct(ctx context.Context, collection *mongo.Collection, id primitive.ObjectID, updateFields bson.M) (int64, error) {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": updateFields}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, err
	}

	return result.ModifiedCount, nil
}

func DeleteProduct(ctx context.Context, collection *mongo.Collection, id primitive.ObjectID) (int64, error) {
	filter := bson.M{"_id": id}

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}

func FetchInventoryByProductID(ctx context.Context, id primitive.ObjectID) (types.InventoryResponse, error) {
	var inventory types.InventoryResponse

	inventoryURL := fmt.Sprintf("http://inventory-service:6000/get/%s", id)

	// ✅ Create an HTTP request with context (timeout-safe)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, inventoryURL, nil)
	// nil means no request body (GET requests usually don’t have a body).
	if err != nil {
		return inventory, fmt.Errorf("failed to create request: %v", err)
	}

	client := http.Client{Timeout: 5 * time.Second}
	// So client := http.Client{Timeout: 5 * time.Second} creates a simple client with a 5-second timeout, meaning the request will fail if it takes longer than 5 seconds.
	resp, err := client.Do(req)
	if err != nil {
		return inventory, fmt.Errorf("inventory service unavailable: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return inventory, fmt.Errorf("inventory not found (status: %d)", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&inventory); err != nil {
		return inventory, fmt.Errorf("failed to decode inventory response: %v", err)
	}

	return inventory, nil
}

type InventoryRequest struct {
	ProductId primitive.ObjectID `json:"productId"`
	Inventory int64              `json:"inventory"`
}

func CreateInventoryForProduct(id primitive.ObjectID) error {
	inventoryURL := "http://inventory-service:6000/create"

	bodyData := InventoryRequest{
		ProductId: id,
		Inventory: 100,
	}

	bodyBytes, err := json.Marshal(bodyData)
	if err != nil {
		return fmt.Errorf("failed to marshal inventory request: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	/*Notice that the body has to be an io.Reader.
	io.Reader is an interface in Go that represents anything you can read bytes from.
	Examples: a file, a network connection, a byte slice, or a buffer.
	So if you have some JSON in a byte slice (like bodyBytes), you need a type that implements io.Reader.
	2️⃣ What bytes.NewBuffer Does
	bytes.NewBuffer converts a []byte into a buffer that implements io.Reader (and io.Writer) so you can pass it to functions that need an io.Reader.*/
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, inventoryURL, bytes.NewBuffer(bodyBytes))

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("inventory service unavailable: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("inventory creation failed with status: %d", resp.StatusCode)
	}

	fmt.Println("✅ Inventory created successfully for product:", id.Hex())
	return nil

}
