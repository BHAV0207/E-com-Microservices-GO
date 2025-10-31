package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Paments struct {
	Id            primitive.ObjectID `bson:"_id, omitempty" json:"_id"`
	OrderID       string             `bson:"orderId" json:"orderId" validate:"required"`
	UserID        string             `bson:"userId" json:"userId" validate:"required"`
	Amount        float64            `bson:"amount" json:"amount" validate:"required"`
	Currency      string             `bson:"currency" json:"currency"`
	Status        string             `bson:"status" json:"status"`
	Method        string             `bson:"method" json:"method"`
	GatewayTxnID  string             `bson:"gatewayTxnId,omitempty" json:"gatewayTxnId,omitempty"`
	FailureReason string             `bson:"failureReason,omitempty" json:"failureReason,omitempty"`
	CreatedAt     time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time          `bson:"updatedAt" json:"updatedAt"`
}

