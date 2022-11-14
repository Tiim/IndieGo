package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"tiim/go-comment-api/api"
	"tiim/go-comment-api/comments"
	"tiim/go-comment-api/event"
	"tiim/go-comment-api/model"
	"tiim/go-comment-api/webmentions"
	"tiim/go-comment-api/wmsend"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	httpClient := &http.Client{Timeout: time.Second * 10}
	scheduler := gocron.NewScheduler(time.UTC)

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
			webmentions.NewTargetChecker("tiim.ch", "localhost"),
			webmentions.NewDomainChecker(wmStore),
			webmentions.NewLinkToTargetChecker(),
			webmentions.NewMicroformatEnricherChecker(),
		},
	)
	wmApi := webmentions.NewApi(wmStore, webmentions.NewMentionsQueueWorker(wmStore, wmChecker))

	scheduler.Every(4).Hours().Do(wmStore.PopulateQueue)

	//
	// Generic Comments
	//

	commentProvider := []api.CommentProvider{
		commentStore,
		wmStore,
	}

	//
	// Sending webmentions
	//

	wmSendStore := wmsend.NewWmSendStore(store.GetDBConnection())
	wmSender := wmsend.NewWmSend(wmSendStore, httpClient, os.Getenv("WM_SEND_RSS_URL"))

	scheduler.Every(1).Hour().Do(wmSender.SendNow)

	//
	// Webhooks
	//

	webhookModule := api.NewWebhookModule()
	webhookModule.RegisterWebhook("page-build", func(c *gin.Context) error {
		wmSender.SendNow()
		return nil
	}, api.NewGithubValidator(os.Getenv("GITHUB_WEBHOOK_SECRET")))

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
		webhookModule,
	}

	log.Println("Starting server")
	scheduler.StartAsync()
	server := api.NewCommentServer(apiModules)
	err = server.Start()
	if err != nil {
		log.Fatal(err)
	}
}
