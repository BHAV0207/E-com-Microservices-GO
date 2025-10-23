package handler

import (
	workerpool "github.com/BHAV0207/user-service/internal/workerPool"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	Collection      *mongo.Collection
	UserCreatedPool *workerpool.WorkerPool
	UserDeletedPool *workerpool.WorkerPool
}
