package api

import (
	"html/template"

	_ "embed"

	"github.com/gin-gonic/gin"
)

//go:embed templates/index.tmpl
var indexTemplate string

type indexModule struct {
	template *template.Template
}

func NewIndexModule() *indexModule {
	template := template.Must(template.New("index").Parse(indexTemplate))
	im := indexModule{template}
	return &im
}

func (ui *indexModule) Name() string {
	return "Index"
}

func (ui *indexModule) Init(r *gin.Engine) error {
	return nil
}

func (ui *indexModule) RegisterRoutes(r *gin.Engine) error {
	r.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		ui.template.Execute(c.Writer, nil)
	})
	r.HEAD("/", func(c *gin.Context) {
		c.Status(200)
	})
	return nil
}
