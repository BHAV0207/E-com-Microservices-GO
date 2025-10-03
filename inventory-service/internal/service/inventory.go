package service

import (
	"context"

	"github.com/BHAV0207/inventory-service/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Get(context context.Context, collection *mongo.Collection, filter bson.M) (models.Inventory, error) {
	var inventory models.Inventory

	err := collection.FindOne(context, filter).Decode(&inventory)
	if err != nil {
		return inventory, err
	}

	return inventory, nil
}
