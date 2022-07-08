package api

import (
	"fmt"
	"log"
	"net/http"
	"tiim/go-comment-api/model"

	"github.com/gin-gonic/gin"
)

type commentServer struct {
	store model.Store
}

func NewCommentServer(store model.Store) *commentServer {
	return &commentServer{store: store}
}

func (cs *commentServer) Start() {
	r := gin.New()
	r.RemoveExtraSlash = true
	r.RedirectTrailingSlash = false

	ui := newAdminRoutes(r, cs.store)

	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(trailingSlash(r))
	r.Use(cors())
	r.GET("/comment", cs.handleGetAllComments)
	r.GET("/comment/:page", cs.handleGetComments)
	r.POST("/comment", cs.handlePostComment)

	ui.start()

	r.Run(":8080")

}

func (cs *commentServer) handlePostComment(c *gin.Context) {
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
	err := cs.store.NewComment(&comment)
	if err != nil {
		fmt.Println("Error inserting comment: ", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, comment)
}

func (cs *commentServer) handleGetAllComments(c *gin.Context) {
	comments, err := cs.store.GetAllComments()
	log.Println("Got all comments: ", comments)
	if err != nil {
		log.Println("Error getting comments: ", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, comments)
}

func (cs *commentServer) handleGetComments(c *gin.Context) {
	uuidParam := c.Param("uuid")
	comments, err := cs.store.GetCommentsForPost(uuidParam)
	if err != nil {
		fmt.Println("Error getting comments for post ", uuidParam, err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, comments)
}
