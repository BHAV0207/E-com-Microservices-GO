package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	UserId    primitive.ObjectID `bson:"userId" json:"userId"`
	Products  []OrderProducts    `bson:"products" json:"products"`
	Address   string             `bson:"address" json:"address"`
	Status    string             `bson:"status" json:"status"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}

type OrderProducts struct {
	ProductId primitive.ObjectID `bson:"productId" json:"productId"`
	Quantity  int64              `bson:"quantity" json:"quantity"`
	Price     float64            `bson:"price" json:"price"`
}
