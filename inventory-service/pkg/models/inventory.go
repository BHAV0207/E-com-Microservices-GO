package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Inventory struct {
	ID        primitive.ObjectID `bsom:"_id,ommitempty" json:"_id"`
	ProductId primitive.ObjectID `bson:"productId" json:"productId"`
	Inventory int64              `bson:"inventory" json:"inventory"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
	CreatedAt time.Time          `bson:"createdAt" json:"ccreatedAt"`
}
