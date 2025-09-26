package service

import (
	"context"

	"github.com/BHAV0207/product-service/pkg/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

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
