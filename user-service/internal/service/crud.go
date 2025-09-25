package service

import (
	"context"

	"github.com/BHAV0207/user-service/pkg/models"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

func GetAll(ctx context.Context, collection *mongo.Collection) ([]models.User, error) {
	/*bson.D{} is an empty filter → matches all documents (like SELECT *).
	bson.D is an ordered slice of key/value pairs; useful when order matters (e.g., for some operators). bson.M is a map (unordered).*/
	cursor, err := collection.Find(ctx, bson.D{})
	/*cursor is a *mongo.Cursor — an iterator over the query result. It does not read all documents into memory at once; it fetches batches from the server as you iterate.*/
	if err != nil {
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

func Update(ctx context.Context , collection *mongo.Collection , id string)(int64 , error){
	
}
