package pkg

import "go.mongodb.org/mongo-driver/bson/primitive"

type CartItem struct {
	ProductId primitive.ObjectID `bson:"productId" json:"productId"`
	Quantity  int64              `bson:"quantity" json:"quantity"`
}

type Cart struct {
	UserId primitive.ObjectID `bson:"userId" json:"userId"`
	Items  []CartItem         `bson:"items" json:"items"`
}
