package service

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/BHAV0207/product-service/pkg/models"
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

func InsertProduct(ctx context.Context, collection *mongo.Collection, product models.Product) (interface{}, error) {
	result, err := collection.InsertOne(ctx, product)
	return result.InsertedID, err
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
