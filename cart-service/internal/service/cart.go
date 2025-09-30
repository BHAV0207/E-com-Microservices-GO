package service

import (
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

func ValidateUser(id string) bool {
	resp, err := http.Get(fmt.Sprintf("http://user-service:8080/%s", id))
	if err != nil || resp.StatusCode != http.StatusOK {
		return false
	}
	return true
}

func ValidateProduct(id string) bool {
	resp, err := http.Get(fmt.Sprintf("http://product-service:4000/get/%s", id))
	if err != nil || resp.StatusCode != http.StatusOK {
		return false
	}
	return true
}

func AddItemToCart(ctx context.Context, collection *mongo.Collection, userId, productId primitive.ObjectID, quantity int64) error {
	var cart struct {
		ID     primitive.ObjectID `bson:"_id,omitempty"`
		UserID primitive.ObjectID `bson:"userId"`
		Items  []struct {
			ProductID primitive.ObjectID `bson:"productId"`
			Quantity  int64              `bson:"quantity"`
		} `bson:"items"`
	}

	err := collection.FindOne(ctx, bson.M{"userId": userId}).Decode(&cart)
	if err == mongo.ErrNoDocuments {
		// 2. No cart exists, create a new one
		newCart := bson.M{
			"userId": userId,
			"items": []bson.M{{
				"productid": productId,
				"quantity":  quantity,
			}},
		}
		_, err := collection.InsertOne(ctx, newCart)
		return err
	} else if err != nil {
		return err
	}

	found := false
	for i, item := range cart.Items {
		if item.ProductID == productId {
			cart.Items[i].Quantity += quantity
			found = true
			break
		}
	}
	if !found {
		cart.Items = append(cart.Items, struct {
			ProductID primitive.ObjectID `bson:"productId"`
			Quantity  int64              `bson:"quantity"`
		}{ProductID: productId, Quantity: quantity})
	}

	// 4. Update the cart in MongoDB
	_, err = collection.UpdateOne(
		ctx,
		bson.M{"user_id": userId},
		bson.M{"$set": bson.M{"items": cart.Items}},
	)
	return err
}
