package handler

import "go.mongodb.org/mongo-driver/mongo"

type PaymentHandler struct{
	Collection *mongo.Collection
}

