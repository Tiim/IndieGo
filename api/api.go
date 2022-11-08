package api

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type apiServer struct {
	modules []ApiModule
}

func NewCommentServer(modules []ApiModule) *apiServer {
	return &apiServer{modules: modules}
}

//go:embed templates/*
var templates embed.FS

func (cs *apiServer) Start() error {
	r := gin.New()
	r.RemoveExtraSlash = true
	r.RedirectTrailingSlash = false

	r.Use(ErrorMiddleware())

	tp := template.Must(template.New("").ParseFS(templates, "templates/*"))
	r.SetHTMLTemplate(tp)

	for _, module := range cs.modules {
		if err := module.Init(r, &templates); err != nil {
			return fmt.Errorf("initialising module %s failed: %w", module.Name(), err)
		}
	}

	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(trailingSlash(r))
	r.Use(cors())

	for _, module := range cs.modules {
		if err := module.RegisterRoutes(r); err != nil {
			return fmt.Errorf("registering routes failed for module %s: %w", module.Name(), err)
		}
	}

	r.Run(":8080")
	return nil
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
