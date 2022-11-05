package api

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"tiim/go-comment-api/model"

	"github.com/gin-gonic/gin"
)

type subscriptionModule struct {
	store model.SubscribtionStore
}

func NewSubscriptionModule(store model.SubscribtionStore) *subscriptionModule {
	sm := subscriptionModule{store: store}
	return &sm
}

func (sm *subscriptionModule) Name() string {
	return "Subscription"
}

func (sm *subscriptionModule) Init(r *gin.Engine, templates fs.FS) error {
	return nil
}

func (sm *subscriptionModule) RegisterRoutes(r *gin.Engine) error {
	r.GET("/unsubscribe/comment/:secret", sm.handleUnsubscribeComment)
	r.GET("/unsubscribe/email/:email", sm.handleUnsubscribeEmail)
	return nil
}

func (sm *subscriptionModule) handleUnsubscribeComment(c *gin.Context) {
	secret := c.Param("secret")

	comment, err := sm.store.Unsubscribe(secret)
	if err != nil {
		fmt.Println("Error unsubscribing comment: ", err)
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("unsubscribing comment failed: %w", err))
		return
	}

	c.HTML(http.StatusOK, "unsubscribe_cmt.tmpl", gin.H{"comment": comment, "emailUrl": template.URLQueryEscaper(comment.Email)})
}

func (sm *subscriptionModule) handleUnsubscribeEmail(c *gin.Context) {
	email := c.Param("email")

	comments, err := sm.store.UnsubscribeAll(email)
	if err != nil {
		fmt.Println("Error unsubscribing comments: ", err)
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("unsubscribing comments failed: %w", err))
		return
	}

	c.HTML(http.StatusOK, "unsubscribe_email.tmpl", gin.H{"comments": comments})
}
