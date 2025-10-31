package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Paments struct {
	Id       primitive.ObjectID `bson:"_id, omitempty" json:"_id"`
	OrderID  string             `bson:"orderId" json:"orderId" validate:"required"`
	UserID   string             `bson:"userId" json:"userId" validate:"required"`
	Amount   float64            `bson:"amount" json:"amount" validate:"required"`
	Currency string             `bson:"currency" json:"currency"`
	Status   string             `bson:"status" json:"status"`
	Method   string             `bson:"method" json:"method"`
	
}
