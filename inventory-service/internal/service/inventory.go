package service

import (
	"context"

	"github.com/BHAV0207/inventory-service/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Get(context context.Context, collection *mongo.Collection, filter bson.M) (models.Inventory, error) {
	var inventory models.Inventory

	err := collection.FindOne(context, filter).Decode(&inventory)
	if err != nil {
		return inventory, err
	}

	return inventory, nil
}

func Create(ctx context.Context, collection *mongo.Collection, inventory models.Inventory) (*mongo.InsertOneResult, error) {
	res, err := collection.InsertOne(ctx, inventory)

	if err != nil {
		return nil, err
	}
	return res, nil
}

func Update(ctx context.Context, collection *mongo.Collection, id primitive.ObjectID, updateFileds map[string]interface{}) (models.Inventory, error) {
	var updatedInventory models.Inventory

	filter := bson.M{"_id": id}

	update := bson.M{
		"$set": updateFileds,
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedInventory)
	if err != nil {
		return models.Inventory{}, err
	}

	return updatedInventory, nil
}
