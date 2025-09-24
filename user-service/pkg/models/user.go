package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name      string             `bson:"name" json:"name" validate:"required,min=2,max=50"`
	Email     string             `bson:"email" json:"email" validate:"required,email"`
	Phone     int                `bson:"phone" json:"phone"`
	Password  string             `bson:"password" json:"-" validate:"required,min=6"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
