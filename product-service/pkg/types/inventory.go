package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type InventoryResponse struct {
    ID        primitive.ObjectID `json:"_id"`
    ProductId primitive.ObjectID `json:"productId"`
    Inventory int64              `json:"inventory"`
    UpdatedAt string             `json:"updatedAt"`
    CreatedAt string             `json:"createdAt"`
}
