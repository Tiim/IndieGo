package main

import (
	"fmt"
	"log"
	"os"
	"tiim/go-comment-api/api"
	"tiim/go-comment-api/comments"
	"tiim/go-comment-api/event"
	"tiim/go-comment-api/model"
	"tiim/go-comment-api/webmentions"

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

	commentToUrl := func(page string, id string) string {
		return fmt.Sprintf("https://tiim.ch/%s#%s", page, id)
	}

	//
	// Comments
	//

	commentStore := comments.NewCommentStore(store.GetDBConnection(), commentToUrl)

	//
	// Webmentions
	//
	wmStore := webmentions.NewStore(store)
	wmChecker := webmentions.NewWebmentionChecker(
		[]webmentions.Checker{
			webmentions.NewDomainChecker(wmStore),
			webmentions.NewLinkToTargetChecker(),
		},
	)
	wmApi := webmentions.NewApi(wmStore, webmentions.NewMentionsQueueWorker(wmStore, wmChecker))

	//
	// Generic Comments
	//

	commentProvider := []api.CommentProvider{
		commentStore,
		wmStore,
	}

	//
	// Event handlers
	//

	emailnotify := &event.EmailNotify{
		From:     os.Getenv("EMAIL_FROM"),
		To:       os.Getenv("EMAIL_NOTIFY_TO"),
		Username: os.Getenv("EMAIL_USERNAME"),
		Password: os.Getenv("EMAIL_PASSWORD"),
		SmtpHost: os.Getenv("EMAIL_SMTP_HOST"),
		SmtpPort: os.Getenv("EMAIL_SMTP_PORT"),
		Subject:  os.Getenv("EMAIL_NOTIFY_SUBJECT"),
	}

	replyEmailNotify := comments.NewReplyEmail(
		commentStore,
		os.Getenv("EMAIL_FROM"),
		os.Getenv("EMAIL_REPLY_SUBJECT"),
		os.Getenv("EMAIL_USERNAME"),
		os.Getenv("EMAIL_PASSWORD"),
		os.Getenv("EMAIL_SMTP_HOST"),
		os.Getenv("EMAIL_SMTP_PORT"),
		os.Getenv("BASE_URL"),
	)

	cleanup := &event.CleanUp{Store: store}

	eventHandler := event.NewHandlerList([]event.Handler{
		emailnotify,
		replyEmailNotify,
		cleanup,
	})

	commentStore.SetEventHandler(eventHandler)
	wmStore.SetEventHandler(eventHandler)

	adminSections := []api.AdminSection{
		comments.NewAdminCommentSection(commentStore),
		api.NewAdminBackupSection(store),
		webmentions.NewAdminWebmentionsSection((wmStore)),
	}

	apiModules := []api.ApiModule{
		api.NewIndexModule(),
		api.NewCommentModule(commentProvider),
		api.NewAdminModule(adminSections),
		comments.NewCommentModule(commentStore),
		comments.NewSubscriptionModule(commentStore),
		wmApi,
	}

	log.Println("Starting server")
	server := api.NewCommentServer(apiModules)
	err = server.Start()
	if err != nil {
		log.Fatal(err)
	}
}
