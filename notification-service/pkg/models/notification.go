package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Notifiction struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	UserId    primitive.ObjectID `bson:"userId" json:"userId"`
	OrderId   primitive.ObjectID `bson:"orderId" json:"orderId"`
	Type      string             `bson:"type" json:"type"`
	Message   string             `bson:"message" json:"message"`
	Status    string             `bson:"status" json:"status"`
	CreatedAt int64              `bson:"createdAt" json:"createdAt"`
	UpdatedAt int64              `bson:"updatedAt" json:"updatedAt"`
}


