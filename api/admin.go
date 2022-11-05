package api

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"tiim/go-comment-api/model"

	"github.com/gin-gonic/gin"
)

type AdminSection interface {
	Init(templates fs.FS) error
	Name() string
	HTML() (string, error)
	RegisterRoutes(r *gin.RouterGroup) error
}

type sectionData struct {
	Name string
	HTML template.HTML
}

type adminModule struct {
	password string
	store    model.Store
	group    *gin.RouterGroup
	sections []AdminSection
}

func NewAdminModule(store model.Store, sections []AdminSection) *adminModule {
	password, envExists := os.LookupEnv("ADMIN_PW")
	if !envExists {
		log.Fatal("env variable ADMIN_PW not found, check .env file")
	}

	uir := adminModule{password: password, store: store, sections: sections}
	return &uir
}

func (ui *adminModule) Name() string {
	return "Admin"
}

func (ui *adminModule) Init(r *gin.Engine, templates fs.FS) error {
	for _, section := range ui.sections {
		if err := section.Init(templates); err != nil {
			return fmt.Errorf("initialising section %s failed: %w", section.Name(), err)
		}
	}
	ui.group = r.Group("/admin", gin.BasicAuth(gin.Accounts{"admin": ui.password}))
	return nil
}

func (ui *adminModule) RegisterRoutes(r *gin.Engine) error {

	for _, section := range ui.sections {
		section.RegisterRoutes(ui.group)
	}

	ui.group.GET("", ui.adminDashboard)
	return nil
}

func (ui *adminModule) adminDashboard(c *gin.Context) {
	sections := make([]sectionData, len(ui.sections))
	for i, section := range ui.sections {
		name := section.Name()
		html, err := section.HTML()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("unable to render section %s: %w", name, err))
			return
		}
		sections[i] = sectionData{Name: name, HTML: template.HTML(html)}
	}

	c.HTML(http.StatusOK, "dashboard.tmpl", gin.H{"sections": sections})
}
