package publicsite

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"tiim/go-comment-api/config"

	"github.com/gin-gonic/gin"
)

//go:embed index.tmpl
var indexTemplate string

//go:embed assets/*
var assets embed.FS

type publicModule struct {
	template        *template.Template
	debugApertureId string
	logger          *log.Logger
}

func newPublicModule(debugApertureId string, logger *log.Logger) *publicModule {
	template := template.Must(template.New("index").Parse(indexTemplate))
	im := publicModule{template: template, debugApertureId: debugApertureId, logger: logger}
	return &im
}

func (ui *publicModule) Name() string {
	return "Index"
}

func (ui *publicModule) Init(config config.GlobalConfig) error {
	return nil
}

func (ui *publicModule) Start() error {
	return nil
}

func (ui *publicModule) RegisterRoutes(r *gin.Engine) error {
	assetsFolder, err := fs.Sub(assets, "assets")
	if err != nil {
		return fmt.Errorf("unable to get assets folder: %w", err)
	}
	r.StaticFS("/assets", http.FS(assetsFolder))

	r.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		ui.template.Execute(c.Writer, gin.H{
			"Debug": gin.H{
				"ApertureId": ui.debugApertureId,
			},
		})
	})
	r.HEAD("/", func(c *gin.Context) {
		c.Status(200)
	})
	return nil
}
