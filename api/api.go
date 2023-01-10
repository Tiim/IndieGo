package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"tiim/go-comment-api/config"

	"github.com/gin-gonic/gin"
)

var logger = log.New(os.Stdout, "[api] ", log.Flags())

type apiServer struct {
	plugins map[string][]config.ModuleInstance
}

func NewApiServer(modules map[string][]config.ModuleInstance) *apiServer {
	return &apiServer{plugins: modules}
}

func (cs *apiServer) Start() (*gin.Engine, error) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.RemoveExtraSlash = true
	r.RedirectTrailingSlash = false

	r.Use(ErrorMiddleware())
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	for name, modules := range cs.plugins {
		for _, module := range modules {
			apiPlugin, ok := module.(config.GroupedApiPluginInstance)
			if !ok {
				continue
			}
			if err := apiPlugin.InitGroups(r); err != nil {
				return nil, fmt.Errorf("initialising module %s failed: %w", name, err)
			}
		}
	}

	r.Use(trailingSlash(r))
	r.Use(cors())

	for name, modules := range cs.plugins {
		for _, module := range modules {
			apiPlugin, ok := module.(config.ApiPluginInstance)
			if !ok {
				continue
			}
			if err := apiPlugin.RegisterRoutes(r); err != nil {
				return nil, fmt.Errorf("registering routes failed for module %s: %w", name, err)
			}
		}
	}

	routes := r.Routes()
	for _, route := range routes {
		logger.Printf("Registered route: %-6s %s", route.Method, route.Path)
	}

	return r, nil
}

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			logger.Printf("Error: %v", c.Errors)
			status := c.Writer.Status()
			if status == 0 || status < 400 {
				status = http.StatusInternalServerError
			}
			c.JSON(status, gin.H{"status": status, "error": c.Errors})
		}
	}
}
