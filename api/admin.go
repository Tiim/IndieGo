package api

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"tiim/go-comment-api/model"
	"time"

	"github.com/gin-gonic/gin"
)

//go:embed templates/*
var templates embed.FS

type uiRoutes struct {
	store model.Store
	group *gin.RouterGroup
}

func newAdminRoutes(r *gin.Engine, store model.Store) uiRoutes {

	password, envExists := os.LookupEnv("ADMIN_PW")
	if !envExists {
		log.Fatal("env variable ADMIN_PW not found, check .env file")
	}

	admin := r.Group("/admin", gin.BasicAuth(gin.Accounts{
		"admin": password,
	}))

	uir := uiRoutes{store: store, group: admin}

	tp := template.Must(template.New("").ParseFS(templates, "templates/*"))
	r.SetHTMLTemplate(tp)

	return uir
}

func (ui *uiRoutes) start() {
	ui.group.GET("", ui.adminDashboard)
	ui.group.POST("delete", ui.deleteComment)
}

func (ui *uiRoutes) adminDashboard(c *gin.Context) {
	comments, err := ui.store.GetAllComments(time.Time{})

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("unable to fetch comments: %w", err))
		return
	}

	c.HTML(http.StatusOK, "dashboard.tmpl", gin.H{"comments": comments})
}

func (ui *uiRoutes) deleteComment(c *gin.Context) {
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
