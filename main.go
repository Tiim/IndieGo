package main

import (
	"log"
	"tiim/go-comment-api/api"
	"tiim/go-comment-api/model"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	store, err := model.NewSQLiteStore()
	if err != nil {
		log.Fatal(err)
	}
	server := api.NewCommentServer(store)
	server.Start()
}
