package handler

import "go.mongodb.org/mongo-driver/mongo"

type UserHandler struct {
    Collection *mongo.Collection
}
