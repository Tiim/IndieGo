package api

import (
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

type indexModule struct {
}

func NewIndexModule() *indexModule {
	im := indexModule{}
	return &im
}

func (ui *indexModule) Name() string {
	return "Index"
}

func (ui *indexModule) Init(r *gin.Engine, templates fs.FS) error {
	return nil
}

func (ui *indexModule) RegisterRoutes(r *gin.Engine) error {
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{})
	})
	return nil
}
