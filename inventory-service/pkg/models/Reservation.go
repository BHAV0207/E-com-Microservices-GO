package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReservationItem struct {
	ProductId string `bson:"productId" json:"productId"`
	Quantity  int64  `bson:"quantity" json:"quantity"`
}

type Reservation struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ReservationID string             `bson:"reservationId" json:"reservationId"`
	OrderID       string             `bson:"orderId" json:"orderId"`
	Status        string             `bson:"status" json:"status"` // PENDING, COMMITTED, CANCELLED
	Items         []ReservationItem  `bson:"items" json:"items"`
	ExpiresAt     time.Time          `bson:"expiresAt" json:"expiresAt"`
	CreatedAt     time.Time          `bson:"createdAt" json:"createdAt"`
}
