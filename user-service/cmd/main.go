package cmd

import (
	"context"
	"log"

	"github.com/BHAV0207/user-service/internal/repository"
)

func main() {
	uri := "mongodb+srv://jainbhav0207_db_user:XB9P4Jgp0fzqBCOS@cluster0.oa5vscu.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"
	client := repository.ConnectDb(uri)
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	// db := client.Database("UserService")

}
