package kafka

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserCreateEvent struct {
	UserID primitive.ObjectID `json:"userId"`
	Email  string             `json:"email"`
}


