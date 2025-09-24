package service

import (
	"context"

	"github.com/BHAV0207/user-service/pkg/models"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

func GetUserByEmail(ctx context.Context, Collection *mongo.Collection, email string) (models.User, bool, error) {
	filter := bson.M{"email": email}
	/*
		bson.M
		bson is a package from the official MongoDB Go driver (go.mongodb.org/mongo-driver/bson).
		M is a type alias for a map[string]interface{} — it’s the most common way to build queries and documents in Go for MongoDB.
		bson.M basically means “a BSON (MongoDB’s data format) document represented as a Go map”.
		👉 So bson.M{"key": value} means a document where "key" has value.
	*/

	var existingUser models.User
	err := Collection.FindOne(ctx, filter).Decode(&existingUser)

	if err == mongo.ErrNoDocuments {
		// ✅ No user found with this email
		return existingUser, false, nil
	}
	if err != nil {
		return existingUser, false, err
	}

	// ✅ User exists
	return existingUser, true, nil
}
