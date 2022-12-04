package admin

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "embed"

	"github.com/gin-gonic/gin"
)

//go:embed dashboard.tmpl
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

type AdminModule struct {
	password string
	group    *gin.RouterGroup
	sections []AdminSection
	template *template.Template
}

func newAdminModule(password string) *AdminModule {
	if password == "" || len(password) < 8 {
		log.Fatal("admin password must be at least 8 characters long")
	}

	template := template.Must(template.New("dashboard").Parse(dashboardTemplate))

	uir := AdminModule{password: password, sections: []AdminSection{}, template: template}
	return &uir
}

func (ui *AdminModule) Name() string {
	return "admin"
}
func (ui *AdminModule) Init() error {
	for _, section := range ui.sections {
		if err := section.Init(); err != nil {
			return fmt.Errorf("initialising section %s failed: %w", section.Name(), err)
		}
	}
	return nil
}

func (ui *AdminModule) InitGroups(r *gin.Engine) error {
	ui.group = r.Group("/admin", gin.BasicAuth(gin.Accounts{"admin": ui.password}))
	return nil
}

func (ui *AdminModule) RegisterRoutes(r *gin.Engine) error {

	for _, section := range ui.sections {
		section.RegisterRoutes(ui.group)
	}

	ui.group.GET("", ui.adminDashboard)
	return nil
}

func (ui *AdminModule) Start() error {
	return nil
}

func (ui *AdminModule) RegisterSection(section AdminSection) {
	ui.sections = append(ui.sections, section)
}

func (ui *AdminModule) adminDashboard(c *gin.Context) {
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
