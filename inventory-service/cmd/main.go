package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/BHAV0207/inventory-service/internal/repository"
)

func main() {
	uri := "mongodb+srv://jainbhav0207_db_user:PdzvcXtnxHW4B3vv@cluster0.exrhn4j.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"

	client := repository.ConnectDb(uri)
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	db := client.Database("OrderService")

	fmt.Println("Server listening on http://localhost:6000")
	log.Fatal(http.ListenAndServe(":6000", router))

}
