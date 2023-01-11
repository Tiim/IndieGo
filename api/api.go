package api

import (
	"fmt"
	"log"
	"os"
	"tiim/go-comment-api/config"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	// prometheus registry
	registry := prometheus.NewRegistry()

	r.Use(errorMiddleware())

	r.GET("/metrics", metricsHandler(registry))
	r.Use(metricsMiddleware(registry))

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

func metricsHandler(reg *prometheus.Registry) gin.HandlerFunc {
	handler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}
