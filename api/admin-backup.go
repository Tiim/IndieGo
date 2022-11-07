package api

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"tiim/go-comment-api/model"
	"time"

	"github.com/gin-gonic/gin"
)

type adminBackupSection struct {
	store    model.BackupStore
	template *template.Template
}

func NewAdminBackupSection(store model.BackupStore) *adminBackupSection {
	return &adminBackupSection{
		store: store,
	}
}

func (ui *adminBackupSection) Init(templates fs.FS) error {
	template := template.Must(template.New("dashboard-backup.tmpl").ParseFS(templates, "templates/dashboard-backup.tmpl"))
	ui.template = template
	return nil
}

func (ui *adminBackupSection) Name() string {
	return "Backup"
}

func (ui *adminBackupSection) HTML() (string, error) {
	var buf bytes.Buffer
	err := ui.template.Execute(&buf, nil)
	if err != nil {
		return "", fmt.Errorf("unable to execute template: %w", err)
	}
	return buf.String(), nil
}

func (ui *adminBackupSection) RegisterRoutes(group *gin.RouterGroup) error {
	group.GET("backup", ui.backup)
	return nil
}

func (ui *adminBackupSection) backup(c *gin.Context) {
	reader, err := ui.store.Backup()
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
