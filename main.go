package main

import (
	"fmt"
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

	err = store.RunMigrations()
	if err != nil {
		log.Fatal(err)
	}

	commentToUrl := func(c model.Comment) string {
		return fmt.Sprintf("https://tiim.ch/%s#%s", c.Page, c.Id)
	}

	emailnotify := &event.EmailNotify{
		From:               os.Getenv("EMAIL_FROM"),
		To:                 os.Getenv("EMAIL_NOTIFY_TO"),
		Username:           os.Getenv("EMAIL_USERNAME"),
		Password:           os.Getenv("EMAIL_PASSWORD"),
		SmtpHost:           os.Getenv("EMAIL_SMTP_HOST"),
		SmtpPort:           os.Getenv("EMAIL_SMTP_PORT"),
		Subject:            os.Getenv("EMAIL_NOTIFY_SUBJECT"),
		CommentToUrlMapper: commentToUrl,
	}

	// TODO: database is locker (SQLITE_BUSY)
	replyEmailNotify := event.NewReplyEmail(
		store,
		os.Getenv("EMAIL_FROM"),
		os.Getenv("EMAIL_REPLY_SUBJECT"),
		os.Getenv("EMAIL_USERNAME"),
		os.Getenv("EMAIL_PASSWORD"),
		os.Getenv("EMAIL_SMTP_HOST"),
		os.Getenv("EMAIL_SMTP_PORT"),
		os.Getenv("BASE_URL"),
		commentToUrl,
	)

	cleanup := &event.CleanUp{Store: store}

	eventStore := event.NewEventStore(store, []event.Handler{
		emailnotify,
		replyEmailNotify,
		cleanup,
	})

	server := api.NewCommentServer(eventStore)
	server.Start()
}
