package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	// UserID      primitive.ObjectID `bson:"userId" json:"userId"`
	Name        string             `bson:"name" json:"name" validate:"required"`
	Description string             `bson:"description" json:"description"`
	Category    string             `bson:"category" json:"category" validate:"required"`
	Price       float64            `bson:"price" json:"price"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
	// Inventory   int64              `bson:"inventory" json:"inventory"`
}
