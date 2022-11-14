package comments

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "embed"

	"github.com/gin-gonic/gin"
)

type subscriptionModule struct {
	store        *commentStore
	unsubComment *template.Template
	unsubEmail   *template.Template
}

//go:embed unsubscribe_cmt.tmpl
var unsubCommentTemplate string

//go:embed unsubscribe_email.tmpl
var unsubEmailTemplate string

func NewSubscriptionModule(store *commentStore) *subscriptionModule {
	unsubComment := template.Must(template.New("unsubComment").Parse(unsubCommentTemplate))
	unsubEmail := template.Must(template.New("unsubEmail").Parse(unsubEmailTemplate))
	sm := subscriptionModule{store: store, unsubComment: unsubComment, unsubEmail: unsubEmail}
	return &sm
}

func (sm *subscriptionModule) Name() string {
	return "Subscription"
}

func (sm *subscriptionModule) Init(r *gin.Engine) error {
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
		log.Println("Error unsubscribing comment: ", err)
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("unsubscribing comment failed: %w", err))
		return
	}

	c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	sm.unsubComment.Execute(c.Writer, map[string]interface{}{"comment": comment, "emailUrl": template.URLQueryEscaper(comment.Email)})
}

func (sm *subscriptionModule) handleUnsubscribeEmail(c *gin.Context) {
	email := c.Param("email")

	comments, err := sm.store.UnsubscribeAll(email)
	if err != nil {
		log.Println("Error unsubscribing comments: ", err)
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("unsubscribing comments failed: %w", err))
		return
	}

	c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	sm.unsubComment.Execute(c.Writer, map[string]interface{}{"comments": comments})
}
