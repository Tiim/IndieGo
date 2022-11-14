package api

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "embed"

	"github.com/gin-gonic/gin"
)

//go:embed templates/dashboard.tmpl
var dashboardTemplate string

type AdminSection interface {
	Init() error
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
	group    *gin.RouterGroup
	sections []AdminSection
	template *template.Template
}

func NewAdminModule(sections []AdminSection) *adminModule {
	password, envExists := os.LookupEnv("ADMIN_PW")
	if !envExists {
		log.Fatal("env variable ADMIN_PW not found, check .env file")
	}

	template := template.Must(template.New("dashboard").Parse(dashboardTemplate))

	uir := adminModule{password: password, sections: sections, template: template}
	return &uir
}

func (ui *adminModule) Name() string {
	return "Admin"
}

func (ui *adminModule) Init(r *gin.Engine) error {
	for _, section := range ui.sections {
		if err := section.Init(); err != nil {
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

	c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	ui.template.Execute(c.Writer, map[string][]sectionData{"sections": sections})
}
