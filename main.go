package main

import (
	"log"
	"tiim/go-comment-api/api"
	"tiim/go-comment-api/model"
)

func main() {
	store, err := model.NewCommentStore()
	if err != nil {
		log.Fatal(err)
	}
	server := api.NewCommentServer(store)
	server.Start()
}
