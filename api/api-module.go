package api

import (
	"io/fs"

	"github.com/gin-gonic/gin"
)

type ApiModule interface {
	Name() string
	Init(r *gin.Engine, templates fs.FS) error
	RegisterRoutes(r *gin.Engine) error
}
