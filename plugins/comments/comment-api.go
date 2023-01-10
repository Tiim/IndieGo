package comments

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"tiim/go-comment-api/config"

	"github.com/gin-gonic/gin"
)

type commentApiModule struct {
	store  commentStore
	logger *log.Logger
}

func NewCommentModule(store commentStore, logger *log.Logger) *commentApiModule {
	im := commentApiModule{store: store, logger: logger}
	return &im
}

func (cm *commentApiModule) Name() string {
	return "comments"
}

func (cm *commentApiModule) Init(config config.GlobalConfig) error {
	return nil
}

func (cm *commentApiModule) InitGroups(r *gin.Engine) error {
	return nil
}

func (cm *commentApiModule) RegisterRoutes(r *gin.Engine) error {
	r.POST("/comment", cm.handlePostComment)
	return nil
}

func (cm *commentApiModule) Start() error {
	return nil
}

func (cm *commentApiModule) handlePostComment(c *gin.Context) {
	var comment comment

	if err := c.BindJSON(&comment); err != nil {
		cm.logger.Println("Error binding comment: ", err)
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("deserialising json failed: %w", err))
		return
	}

	if comment.Content == "" || comment.Page == "" {
		cm.logger.Println("Content or Page is empty")
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("content or page is empty"))
		return
	}

	if len(comment.Content) > 1024 || len(comment.Page) > 50 || len(comment.Name) > 70 || len(comment.Email) > 60 || len(comment.ReplyTo) > 40 {
		cm.logger.Println("Content, Page, Name or Email is too long")
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("content, page, name or email is too long"))
		return
	}

	comment.Page = strings.TrimPrefix(comment.Page, "/")

	err := cm.store.NewComment(&comment)
	if err != nil {
		cm.logger.Println("Error inserting comment: ", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, comment)
}
