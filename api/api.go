package api

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"tiim/go-comment-api/config"

	"github.com/gin-gonic/gin"
)

//go:embed assets/*
var assets embed.FS

type apiServer struct {
	plugins []config.PluginInstance
}

func NewApiServer(modules []config.PluginInstance) *apiServer {
	return &apiServer{plugins: modules}
}

func (cs *apiServer) Start() (*gin.Engine, error) {
	r := gin.New()
	r.RemoveExtraSlash = true
	r.RedirectTrailingSlash = false

	assetsFolder, err := fs.Sub(assets, "assets")
	if err != nil {
		return nil, fmt.Errorf("unable to get assets folder: %w", err)
	}
	r.StaticFS("/assets", http.FS(assetsFolder))

	r.Use(ErrorMiddleware())
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	for _, module := range cs.plugins {
		apiPlugin, ok := module.(config.GroupedApiPluginInstance)
		if !ok {
			continue
		}
		if err := apiPlugin.InitGroups(r); err != nil {
			return nil, fmt.Errorf("initialising module %s failed: %w", module.Name(), err)
		}
	}

	r.Use(trailingSlash(r))
	r.Use(cors())

	for _, module := range cs.plugins {
		apiPlugin, ok := module.(config.ApiPluginInstance)
		if !ok {
			continue
		}
		if err := apiPlugin.RegisterRoutes(r); err != nil {
			return nil, fmt.Errorf("registering routes failed for module %s: %w", module.Name(), err)
		}
	}

	return r, nil
}

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			log.Printf("Error: %v", c.Errors)
			status := c.Writer.Status()
			if status == 0 || status < 400 {
				status = http.StatusInternalServerError
			}
			c.JSON(status, gin.H{"status": status, "error": c.Errors})
		}
	}
}
