package api

import (
	"fmt"
	"log"
	"net/http"
	"tiim/go-comment-api/model"
	"time"

	"github.com/gin-gonic/gin"
)

type commentModule struct {
	store model.Store
}

func NewCommentModule(store model.Store) *commentModule {
	im := commentModule{store: store}
	return &im
}

func (cm *commentModule) Name() string {
	return "Comment"
}

func (cm *commentModule) Init(r *gin.Engine) error {
	return nil
}

func (cm *commentModule) RegisterRoutes(r *gin.Engine) error {
	r.GET("/comment", cm.handleGetAllComments)
	r.GET("/comment/*page", cm.handleGetComments)
	r.POST("/comment", cm.handlePostComment)
	return nil
}

func (cm *commentModule) handlePostComment(c *gin.Context) {
	var comment model.Comment

	if err := c.BindJSON(&comment); err != nil {
		log.Println("Error binding comment: ", err)
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("deserialising json failed: %w", err))
		return
	}

	if comment.Content == "" || comment.Page == "" {
		fmt.Println("Content or Page is empty")
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("content or page is empty"))
		return
	}

	if len(comment.Content) > 1024 || len(comment.Page) > 50 || len(comment.Name) > 70 || len(comment.Email) > 60 || len(comment.ReplyTo) > 40 {
		fmt.Println("Content, Page, Name or Email is too long")
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("content, page, name or email is too long"))
		return
	}
	err := cm.store.NewComment(&comment)
	if err != nil {
		fmt.Println("Error inserting comment: ", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, comment)
}

func (cm *commentModule) handleGetAllComments(c *gin.Context) {
	sinceStr := c.Query("since")
	var since time.Time
	if sinceStr == "" {
		since = time.Time{}
	} else {
		var err error
		since, err = time.Parse(time.RFC3339, sinceStr)
		if err != nil {
			fmt.Println("Error parsing since: ", err)
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid since query parameter: %w", err))
			return
		}
	}

	comments, err := cm.store.GetAllComments(since)
	if err != nil {
		log.Println("Error getting comments: ", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, comments)
}

func (cm *commentModule) handleGetComments(c *gin.Context) {
	page := c.Param("page")
	if page[0] == '/' {
		page = page[1:]
	}
	fmt.Println("uuidParam: ", page)

	sinceStr := c.Query("since")
	var since time.Time
	if sinceStr == "" {
		since = time.Time{}
	} else {
		var err error
		since, err = time.Parse(time.RFC3339, sinceStr)
		if err != nil {
			fmt.Println("Error parsing since: ", err)
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid since query parameter: %w", err))
			return
		}
	}

	comments, err := cm.store.GetCommentsForPost(page, since)
	if err != nil {
		fmt.Println("Error getting comments for post ", page, err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, comments)
}
