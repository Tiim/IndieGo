package api

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"tiim/go-comment-api/model"
	"time"

	"github.com/gin-gonic/gin"
)

type adminModule struct {
	password string
	store    model.Store
	group    *gin.RouterGroup
}

func NewAdminModule(store model.Store) *adminModule {
	password, envExists := os.LookupEnv("ADMIN_PW")
	if !envExists {
		log.Fatal("env variable ADMIN_PW not found, check .env file")
	}

	uir := adminModule{password: password, store: store}
	return &uir
}

func (ui *adminModule) Name() string {
	return "Admin"
}

func (ui *adminModule) Init(r *gin.Engine) error {
	ui.group = r.Group("/admin", gin.BasicAuth(gin.Accounts{"admin": ui.password}))
	return nil
}

func (ui *adminModule) RegisterRoutes(r *gin.Engine) error {
	ui.group.GET("", ui.adminDashboard)
	ui.group.POST("delete", ui.deleteComment)
	ui.group.GET("backup", ui.backup)
	return nil
}

func (ui *adminModule) adminDashboard(c *gin.Context) {
	comments, err := ui.store.GetAllComments(time.Time{})

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("unable to fetch comments: %w", err))
		return
	}

	c.HTML(http.StatusOK, "dashboard.tmpl", gin.H{"comments": comments})
}

func (ui *adminModule) deleteComment(c *gin.Context) {
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

func (ui *adminModule) backup(c *gin.Context) {
	backupStore, ok := (ui.store).(model.BackupStore)
	if !ok {
		log.Printf("Store is not a BackupStore")
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("store is not a backup store"))
		return
	}
	reader, err := backupStore.Backup()
	if err != nil {
		log.Printf("Error backing up: %v", err)
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("unable to backup: %w", err))
		return
	}
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Printf("Error reading backup: %v", err)
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("unable to read backup: %w", err))
		return
	}
	time := time.Now().UTC().Format(time.RFC3339)
	c.Header("Content-Disposition", "attachment; filename=comment-api-"+time+".sqlite")
	c.Header("Content-Length", fmt.Sprintf("%d", len(bytes)))
	c.Data(http.StatusOK, "application/x-sqlite3", bytes)
}
