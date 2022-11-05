package api

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"tiim/go-comment-api/model"
	"time"

	"github.com/gin-gonic/gin"
)

type adminCommentsSection struct {
	store    model.Store
	template *template.Template
}

func NewAdminCommentSection(store model.Store) *adminCommentsSection {
	return &adminCommentsSection{
		store: store,
	}
}

func (ui *adminCommentsSection) Init(templates fs.FS) error {
	template, err := template.New("dashboard-comments.tmpl").ParseFS(templates, "templates/dashboard-comments.tmpl")
	if err != nil {
		return fmt.Errorf("unable to parse template: %w", err)
	}
	ui.template = template
	return nil
}

func (ui *adminCommentsSection) Name() string {
	return "Comments"
}

func (ui *adminCommentsSection) HTML() (string, error) {
	comments, err := ui.store.GetAllComments(time.Time{})
	if err != nil {
		return "", fmt.Errorf("unable to get comments: %w", err)
	}
	var buf bytes.Buffer
	err = ui.template.Execute(&buf, comments)
	if err != nil {
		return "", fmt.Errorf("unable to execute template: %w", err)
	}
	return buf.String(), nil
}

func (ui *adminCommentsSection) RegisterRoutes(group *gin.RouterGroup) error {
	group.POST("delete", ui.deleteComment)
	return nil
}

func (ui *adminCommentsSection) deleteComment(c *gin.Context) {
	commentId := c.PostForm("commentId")

	if commentId == "" {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("no commentId field"))
		return
	}

	if err := ui.store.DeleteComment(commentId); err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("unable to delete comment: %w", err))
		return
	}

	c.Redirect(http.StatusFound, "/admin")
}
