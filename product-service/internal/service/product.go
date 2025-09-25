package service

import (
	"context"

	"github.com/BHAV0207/product-service/pkg/models"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertProduct(ctx context.Context, collection *mongo.Collection, product models.Product) (interface{}, error) {
	result, err := collection.InsertOne(ctx, product)
	return result.InsertedID, err
}
