package comments

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"time"

	_ "embed"

	"github.com/gin-gonic/gin"
)

type adminCommentsSection struct {
	store    *commentStore
	template *template.Template
}

func NewAdminCommentSection(store *commentStore) *adminCommentsSection {
	return &adminCommentsSection{
		store: store,
	}
}

//go:embed admin-comments-section.tmpl
var commentsTemplate string

func (ui *adminCommentsSection) Init(templates fs.FS) error {
	template := template.Must(template.New("comments").Parse(commentsTemplate))
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
