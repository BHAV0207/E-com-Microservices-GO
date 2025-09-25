package main

import (
	"context"
	"log"

	"github.com/BHAV0207/product-service/internal/repository"
)

func main() {
	uri := "mongodb+srv://jainbhav0207_db_user:e94JtmF1QEGw0pxT@cluster0.uxo15jw.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"

	client := repository.ConnectDb(uri)
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	db := client.Database("ProductService")

	// productHandler := &
}
