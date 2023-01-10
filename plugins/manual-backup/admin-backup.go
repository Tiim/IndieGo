package manualbackup

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"tiim/go-comment-api/model"
	"time"

	_ "embed"

	"github.com/gin-gonic/gin"
)

//go:embed dashboard-backup.tmpl
var backupTemplate string

type adminBackupSection struct {
	store    model.BackupStore
	template *template.Template
	logger   *log.Logger
}

func newAdminBackupSection(store model.BackupStore, logger *log.Logger) *adminBackupSection {
	return &adminBackupSection{
		store:  store,
		logger: logger,
	}
}

func (ui *adminBackupSection) Init() error {
	template := template.Must(template.New("dashboard-backup").Parse(backupTemplate))
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
		ui.logger.Printf("Error backing up: %v", err)
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("unable to backup: %w", err))
		return
	}
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		ui.logger.Printf("Error reading backup: %v", err)
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("unable to read backup: %w", err))
		return
	}
	time := time.Now().UTC().Format(time.RFC3339)
	c.Header("Content-Disposition", "attachment; filename=comment-api-"+time+".sqlite")
	c.Header("Content-Length", fmt.Sprintf("%d", len(bytes)))
	c.Data(http.StatusOK, "application/x-sqlite3", bytes)
}
