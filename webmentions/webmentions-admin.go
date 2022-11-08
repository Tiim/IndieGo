package webmentions

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"

	_ "embed"
)

//go:embed admin-webmentions-section.tmpl
var webmentionsTemplate string

type adminWebmentionsSection struct {
	store    *webmentionsStore
	template *template.Template
}

func NewAdminWebmentionsSection(store *webmentionsStore) *adminWebmentionsSection {
	return &adminWebmentionsSection{
		store: store,
	}
}

func (ui *adminWebmentionsSection) Init(templates fs.FS) error {
	template := template.Must(template.New("webmentions").Parse(webmentionsTemplate))
	ui.template = template
	return nil
}

func (ui *adminWebmentionsSection) Name() string {
	return "Webmentions"
}

func (ui *adminWebmentionsSection) HTML() (string, error) {
	wms, err := ui.store.GetWebmentions()
	if err != nil {
		return "", fmt.Errorf("unable to get webmentions: %w", err)
	}

	denylist, err := ui.store.GetDomainDenyList()
	if err != nil {
		return "", fmt.Errorf("unable to get domain deny list: %w", err)
	}

	var buf bytes.Buffer
	err = ui.template.Execute(&buf, map[string]interface{}{"Webmentions": wms, "DenyList": denylist})
	if err != nil {
		return "", fmt.Errorf("unable to execute template: %w", err)
	}
	return buf.String(), nil
}

func (ui *adminWebmentionsSection) RegisterRoutes(group *gin.RouterGroup) error {
	group.POST("/wm/delete", ui.handleDeleteWebmention)
	group.POST("/wm/denylist", ui.handleDenyListWebmention)
	group.POST("/wm/denylist-remove", ui.handleDenyListRemoveDomain)
	return nil
}

func (ui *adminWebmentionsSection) handleDeleteWebmention(c *gin.Context) {
	id := c.PostForm("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "missing id"})
		return
	}

	err := ui.store.DeleteWebmention(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/admin")
}

func (ui *adminWebmentionsSection) handleDenyListWebmention(c *gin.Context) {
	id := c.PostForm("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "missing id"})
		return
	}

	wm, err := ui.store.GetWebmention(id, nil)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	err = ui.store.DenyListDomain(wm.SourceUrl().Hostname())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/admin")
}

func (ui *adminWebmentionsSection) handleDenyListRemoveDomain(c *gin.Context) {
	domain := c.PostForm("domain")
	if domain == "" {
		c.JSON(400, gin.H{"error": "missing domain"})
		return
	}

	err := ui.store.DeleteDomainFromDenyList(domain)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/admin")
}
