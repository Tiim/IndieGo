package main

import (
	"log"
	"os"
	"tiim/go-comment-api/api"
	"tiim/go-comment-api/event"
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

	emailnotify := &event.EmailNotify{
		From:     os.Getenv("EMAIL_FROM"),
		To:       os.Getenv("EMAIL_NOTIFY_TO"),
		Password: os.Getenv("EMAIL_PASSWORD"),
		SmtpHost: os.Getenv("EMAIL_SMTP_HOST"),
		SmtpPort: os.Getenv("EMAIL_SMTP_PORT"),
		Subject:  "[Website] New Comment",
	}

	eventStore := event.NewEventStore(store, []event.Handler{emailnotify})

	server := api.NewCommentServer(eventStore)
	server.Start()
}
