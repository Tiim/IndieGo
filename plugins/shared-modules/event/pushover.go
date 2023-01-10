package event

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"tiim/go-comment-api/model"
)

type pushoverNotify struct {
	apiToken string
	userKey  string
	client   http.Client
	logger   *log.Logger
}

func newPushoverNotify(apiToken, userKey string, client http.Client, logger *log.Logger) *pushoverNotify {
	return &pushoverNotify{
		apiToken: apiToken,
		userKey:  userKey,
		client:   client,
		logger:   logger,
	}
}

func (n *pushoverNotify) OnNewComment(c *model.GenericComment) (bool, error) {
	if n.apiToken != "" && n.userKey != "" {
		go n.doSendNotification(c)
	}
	return true, nil
}

func (n *pushoverNotify) doSendNotification(c *model.GenericComment) {

	title := "New comment"
	message := fmt.Sprintf("New comment on %s by %s<%s>:\n%s\n%s", c.Page, c.Name, c.FromEmail, c.Url, c.Content)

	n.client.PostForm("https://api.pushover.net/1/messages.json", url.Values{
		"token":   {n.apiToken},
		"user":    {n.userKey},
		"message": {message},
		"title":   {title},
	})
}

func (n *pushoverNotify) OnDeleteComment(c *model.GenericComment) (bool, error) {
	return true, nil
}

func (n *pushoverNotify) Name() string {
	return "PushoverNotify"
}
