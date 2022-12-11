package commentprovider

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"tiim/go-comment-api/config"
	"tiim/go-comment-api/model"
	"time"

	"github.com/gin-gonic/gin"
)

type genericCommentApiModule struct {
	CommentProviders []CommentProvider
}

func newCommentProviderModule(CommentProviders []CommentProvider) *genericCommentApiModule {
	return &genericCommentApiModule{CommentProviders: CommentProviders}
}

func (cm *genericCommentApiModule) Name() string {
	return "Comment"
}

func (cm *genericCommentApiModule) Init(c config.GlobalConfig) error {
	return nil
}

func (cm *genericCommentApiModule) Start() error {
	return nil
}

func (cm *genericCommentApiModule) RegisterRoutes(r *gin.Engine) error {
	r.GET("/comment", cm.handleGetAllComments)
	r.GET("/comment/*page", cm.handleGetComments)
	return nil
}

func (cm *genericCommentApiModule) handleGetAllComments(c *gin.Context) {
	sinceStr := c.Query("since")
	var since time.Time
	if sinceStr == "" {
		since = time.Time{}
	} else {
		var err error
		since, err = time.Parse(time.RFC3339, sinceStr)
		if err != nil {
			log.Println("Error parsing since: ", err)
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid since query parameter: %w", err))
			return
		}
	}

	allComments := make([]model.GenericComment, 0)

	for _, commentProvider := range cm.CommentProviders {
		comments, err := commentProvider.GetAllGenericComments(since)
		if err != nil {
			log.Println("Error getting comments: ", err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		allComments = append(allComments, comments...)
	}
	sort.Slice(allComments, func(i, j int) bool {
		return allComments[i].Timestamp > allComments[j].Timestamp
	})
	c.JSON(http.StatusOK, allComments)
}

func (cm *genericCommentApiModule) handleGetComments(c *gin.Context) {
	page := c.Param("page")
	if page[0] == '/' {
		page = page[1:]
	}
	sinceStr := c.Query("since")
	var since time.Time
	if sinceStr == "" {
		since = time.Time{}
	} else {
		var err error
		since, err = time.Parse(time.RFC3339, sinceStr)
		if err != nil {
			log.Println("Error parsing since: ", err)
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid since query parameter: %w", err))
			return
		}
	}

	allComments := make([]model.GenericComment, 0)

	for _, commentProvider := range cm.CommentProviders {
		comments, err := commentProvider.GetGenericCommentsForPage(page, since)
		if err != nil {
			log.Println("Error getting comments: ", err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		allComments = append(allComments, comments...)
	}
	sort.Slice(allComments, func(i, j int) bool {
		return allComments[i].Timestamp > allComments[j].Timestamp
	})
	c.JSON(http.StatusOK, allComments)
}
