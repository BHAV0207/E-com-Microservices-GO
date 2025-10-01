package service

import (
	"context"
	"fmt"

	"github.com/BHAV0207/user-service/pkg/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

func GetAll(ctx context.Context, collection *mongo.Collection) ([]models.User, error) {
	/*bson.D{} is an empty filter → matches all documents (like SELECT *).
	bson.D is an ordered slice of key/value pairs; useful when order matters (e.g., for some operators). bson.M is a map (unordered).*/
	/*cursor is a *mongo.Cursor — an iterator over the query result. It does not read all documents into memory at once; it fetches batches from the server as you iterate.*/
	fmt.Println("📡 Calling Find() on:", collection.Name())
	cursor, err := collection.Find(ctx, bson.M{})

	if err != nil {
		fmt.Println("❌ MongoDB Find error:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func Update(ctx context.Context, collection *mongo.Collection, id primitive.ObjectID, updateFields bson.M) (int64, error) {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": updateFields}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, err
	}

	return result.ModifiedCount, nil
}

func Delete(ctx context.Context, collection *mongo.Collection, id primitive.ObjectID) (int64, error) {
	filter := bson.M{"_id": id}

	res, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, nil

}

func GetById(ctx context.Context, collection *mongo.Collection, id primitive.ObjectID) (models.User, error) {
	var user models.User
	filter := bson.M{"_id": id}

	if err := collection.FindOne(ctx, filter).Decode(&user); err != nil {
		return user, err
	}

	return user, nil
}
